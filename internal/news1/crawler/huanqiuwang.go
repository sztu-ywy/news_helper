package crawler

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/spf13/cast"
	"github.com/uozi-tech/cosy/logger"
	"net/http"
	"news_helper/model"
	"strings"
	"time"
)

var CountryKeywords = []string{"美国", "中国"}
var ScienceKeywords = []string{"人工智能", "科技", "芯片", "deepseek", "kimi", "qwen", "agent", "chatgpt"}
var EconomyKeywords = []string{"军事", "人工智能", "国际", "科技", "美国", "中国", "航母"}
var CompanyKeywords = []string{"华为", "阿里巴巴", "百度", "字节跳动", "英伟达", "特斯拉", "苹果", "谷歌", "微软", "亚马逊", "深度求索", "openai", "chatgpt", "TikTok"}

func HuanQiuCrawer(url, prefix string) []model.News {
	var newsList []model.News
	body, err := CrawlHtml("GET", url, nil, nil, nil)
	if err != nil {
		logger.Errorf("请求失败:", err)
	}

	// 2. 使用 goquery 解析 HTML
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		logger.Errorf("解析 HTML 失败:", err)
	}

	// 3. 查找新闻链接（通常在 a 标签中）
	found := 0
	doc.Find("a[href]").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		text := strings.TrimSpace(s.Text())
		// 过滤无效链接
		if href == "" || len(text) < 5 {
			return
		}
		// href 就是新闻的链接，可能是绝对路径，也可能是相对路径
		// 检查是否为新闻文章链接（可进一步优化）
		if !strings.HasPrefix(href, "/") && !strings.Contains(href, prefix) {
			return
		}
		if strings.Contains(href, "video") || strings.Contains(href, "gallery") {
			return // 排除视频图集
		}
		/*
			有几种结构：
			1. https://capital.huanqiu.com/article/4Nu16NMI9cr
			2. //capital.huanqiu.com/article/4Nu16NMI9cr
		*/

		// 检查标题是否包含任一关键词
		if containsAnyKeyword(text, CountryKeywords) || containsAnyKeyword(text, ScienceKeywords) || containsAnyKeyword(text, EconomyKeywords) || containsAnyKeyword(text, CompanyKeywords) {
			// 补全相对链接
			if strings.HasPrefix(href, "//") {
				// 说明href是相对路径，需要补全环球网的前缀
				href = "https:" + href
			} else if !strings.HasPrefix(href, "http") {
				href = "https://" + href
			}
			// logger.Debugf("✅ 匹配新闻 [%d]:\n", found+1)
			// logger.Debugf("标题: %s\n", text)
			// logger.Debugf("链接: %s\n", href)

			news := &model.News{Title: text, Link: href}

			if strings.HasSuffix(news.Link, prefix) {
				return
			}

			news, err = FetchHuanQiuArticle(news.Link)
			if err != nil {
				logger.Errorf("获取新闻失败: %s", err)
				return
			}
			newsList = append(newsList, *news)
			found++
		}
	})
	logger.Debugf("共找到 %d 条匹配的新闻。\n", found)
	for _, news := range newsList {
		logger.Debugf("新闻链接: %s", news.Link)
		logger.Debugf("新闻标题: %s", news.Title)
		logger.Debugf("新闻时间: %s", news.Time)
	}
	return newsList
}

// 检查文本是否包含任意一个关键词
func containsAnyKeyword(text string, keywords []string) bool {
	for _, kw := range keywords {
		if strings.Contains(text, kw) {
			return true
		}
	}
	return false
}

// FetchHuanQiuArticle 根据新闻链接提取内容
func FetchHuanQiuArticle(articleURL string) (*model.News, error) {
	time.Sleep(100 * time.Millisecond)
	// 1. 发起 HTTP 请求
	resp, err := http.Get(articleURL)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP 错误: %d", resp.StatusCode)
	}

	// 2. 使用 goquery 解析 HTML
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("解析 HTML 失败: %w", err)
	}

	news := &model.News{Link: articleURL}

	// 提取函数：获取指定 class 的 textarea 内容
	extract := func(class string) string {
		var text string
		doc.Find("textarea." + class).Each(func(i int, s *goquery.Selection) {
			text = strings.TrimSpace(s.Text())
		})
		return text
	}

	// 3. 逐个提取字段
	news.Title = extract("article-title")
	news.Content = extract("article-content")
	news.Time = cast.ToUint64(extract("article-ext-xtime"))

	// 处理来源（可能包含 HTML 标签）
	sourceHTML := extract("article-source-name")
	if sourceHTML != "" {
		// 使用 goquery 提取链接文本
		sourceDoc, err := goquery.NewDocumentFromReader(strings.NewReader(sourceHTML))
		if err == nil {
			news.Source = strings.TrimSpace(sourceDoc.Find("a").Text())
		} else {
			// 备用：简单提取
			news.Source = strings.TrimPrefix(sourceHTML, "<a href=\"")
			if idx := strings.Index(news.Source, ">"); idx != -1 {
				news.Source = news.Source[idx+1:]
			}
			if idx := strings.LastIndex(news.Source, "<"); idx != -1 {
				news.Source = news.Source[:idx]
			}
		}
	}

	// 5. 验证是否成功提取标题（基本字段）
	// if news.Title == "" {
	// 	return nil, fmt.Errorf("未找到新闻标题，可能页面结构已变更")
	// }

	return news, nil
}
