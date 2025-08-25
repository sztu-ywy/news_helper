package news1

import (
	"context"
	"encoding/json"
	"fmt"
	"news_helper/internal/news1/crawler"
	"news_helper/internal/queue"
	"news_helper/internal/smtp"

	"news_helper/model"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/robfig/cron/v3"

	"github.com/spf13/cast"
	"github.com/uozi-tech/cosy/logger"
	"github.com/uozi-tech/cosy/redis"
)

const (
	CrawlerRedisPrefix = "crawler"
	CrawlerExpiredTime = 24 * 30 * time.Hour
	CrawlerUser        = "crawler:user:"
	CrawlerTask        = "crawler:task:"
)

// EmailTask 邮件任务结构
type EmailTask struct {
	UserID   uint64 `json:"user_id"`
	TaskID   uint64 `json:"task_id"`
	Email    string `json:"email"`
	Subject  string `json:"subject"`
	Content  string `json:"content"`
	NewsData string `json:"news_data"` // 新闻数据JSON
}

var (
	TaskQueue        *queue.Queue[model.Task]
	NewsCrawlerQueue *queue.Queue[model.News]
	EmailQueue       *queue.Queue[EmailTask] // 新增邮件队列
)

func InitQueue(ctx context.Context) {
	logger.Debug("初始化所有队列")
	NewsCrawlerQueue = queue.New[model.News]("news_crawler", queue.LeftToRight)
	TaskQueue = queue.New[model.Task]("task_queue", queue.LeftToRight)
	EmailQueue = queue.New[EmailTask]("email_queue", queue.LeftToRight) // 初始化邮件队列

	go ConsumeTaskQueue(ctx)
	go ConsumeEmailQueue(ctx) // 启动邮件队列消费者
}
func ConsumeTaskQueue(ctx context.Context) {

	ch, err := TaskQueue.Subscribe(ctx)
	if err != nil {
		logger.Errorf("consume task queue error: %v", err)
		return
	}
	for {
		select {
		case <-ctx.Done():
			return
		case task := <-ch:
			logger.Infof("task: %v", task)
			// todo
			func() {
				defer func() {
					if err := recover(); err != nil {
						buf := make([]byte, 1024)
						runtime.Stack(buf, false)
						logger.Errorf("consume task queue error: %v,buf: %s", err, buf)
					}
				}()
				str_taskId := cast.ToString(task.ID)
				str_userId := cast.ToString(task.UserID)
				a, err := redis.HSet(CrawlerRedisPrefix, str_taskId, str_userId)
				if err != nil {
					logger.Errorf("1HSet error: %v,a: %v", err, a)
				}
				jsonData, err := json.Marshal(task.MediaIDs)
				if err != nil {
					logger.Errorf("json marshal error: %v", err)
				}
				a, err = redis.HSet(CrawlerUser+str_userId, str_taskId, jsonData)
				if err != nil {
					logger.Errorf("2HSet error: %v,a: %v", err, a)
				}
				err = redis.Set(CrawlerTask+str_taskId, task.DailyTime, CrawlerExpiredTime)
				if err != nil {
					logger.Errorf("3HSet error: %v,a: %v", err, a)
				}
				TaskTimer(ctx, str_userId, []string{str_taskId})
			}()
		}
	}
}

func TaskTimer(ctx context.Context, userId string, taskIds []string) {

	// todo
	// tasksTime := make([]string, 0)
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		logger.Error("加载时区失败:", err)
		loc = time.Local
	}
	cronInstance := cron.New(cron.WithLocation(loc))

	logger.Info("doTask")

	for _, taskId := range taskIds {
		taskTime, err := redis.Get(CrawlerTask + taskId)
		if err != nil {
			logger.Error("Get", err)
			continue
		}
		logger.Debug("taskTime:", taskTime, "taskId", taskId)
		RegisterTask(ctx, cronInstance, userId, taskId, taskTime)
	}

	// 启动 cron 定时器
	cronInstance.Start()
	logger.Infof("定时任务已启动，用户 %s 共注册 %d 个任务", userId, len(taskIds))

	// 添加一个测试任务，每分钟执行一次，用于验证定时器是否正常工作
	testEntryID, err := cronInstance.AddFunc("* * * * *", func() {
		logger.Infof("测试定时任务执行 - 用户: %s, 时间: %s", userId, time.Now().Format("2006-01-02 15:04:05"))
	})
	if err != nil {
		logger.Errorf("添加测试定时任务失败: %v", err)
	} else {
		logger.Infof("测试定时任务已添加，EntryID: %d", testEntryID)
	}

	// 监听上下文取消，优雅关闭定时器
	go func() {
		<-ctx.Done()
		logger.Infof("停止用户 %s 的定时任务", userId)
		cronInstance.Stop()
	}()
}

func RegisterTask(ctx context.Context, c *cron.Cron, userId, taskId string, taskTime string) {
	cronTime, err := toCronExpression(taskTime)
	if err != nil {
		logger.Errorf("转换 cron 表达式失败: %v, taskTime: %s", err, taskTime)
		return
	}

	logger.Debugf("注册定时任务 - UserID: %s, TaskID: %s, CronTime: %s, TaskTime: %s",
		userId, taskId, cronTime, taskTime)

	entryID, err := c.AddFunc(cronTime, func() {
		logger.Infof("定时任务触发 - UserID: %s, TaskID: %s, 时间: %s", userId, taskId, time.Now().Format("2006-01-02 15:04:05"))
		doTask(userId, taskId)
	})

	if err != nil {
		logger.Errorf("添加定时任务失败: %v, cronTime: %s", err, cronTime)
	} else {
		logger.Infof("定时任务注册成功 - UserID: %s, TaskID: %s, EntryID: %d, 执行时间: %s",
			userId, taskId, entryID, taskTime)
	}
}

// ConsumeEmailQueue 邮件队列消费者
func ConsumeEmailQueue(ctx context.Context) {
	ch, err := EmailQueue.Subscribe(ctx)
	if err != nil {
		logger.Errorf("consume email queue error: %v", err)
		return
	}
	for {
		select {
		case <-ctx.Done():
			return
		case emailTask := <-ch:
			logger.Infof("处理邮件任务: %+v", emailTask)
			func() {
				defer func() {
					if err := recover(); err != nil {
						buf := make([]byte, 1024)
						runtime.Stack(buf, false)
						logger.Errorf("邮件队列处理错误: %v, stack: %s", err, buf)
					}
				}()
				// 发送邮件的具体实现
				err := sendEmail(emailTask)
				if err != nil {
					logger.Errorf("发送邮件失败: %v", err)
				} else {
					logger.Infof("邮件发送成功: UserID=%d, TaskID=%d", emailTask.UserID, emailTask.TaskID)
				}
			}()
		}
	}
}

func doTask(userId, taskId string) {
	logger.Infof("开始执行定时任务: UserID=%s, TaskID=%s", userId, taskId)

	// 1. 获取任务相关的媒体列表
	mediaListStr, err := redis.HGet(CrawlerUser+userId, taskId)
	if err != nil {
		logger.Errorf("获取媒体列表失败: %v", err)
		return
	}

	var mediaIDs []uint64
	if err := json.Unmarshal([]byte(mediaListStr), &mediaIDs); err != nil {
		logger.Errorf("解析媒体ID列表失败: %v", err)
		return
	}

	if len(mediaIDs) == 0 {
		logger.Warn("媒体列表为空，跳过任务")
		return
	}

	// 2. 理由 id 从数据库中获取媒体信息
	allNewsData := make([]*model.News, 0)
	for _, mediaID := range mediaIDs {
		newsData, err := getMediaInfo(mediaID)
		if err != nil {
			logger.Errorf("爬取媒体 %d 新闻失败: %v", mediaID, err)
			continue
		}
		allNewsData = append(allNewsData, &newsData)
	}

	if len(allNewsData) == 0 {
		logger.Warn("没有爬取到任何新闻数据")
		return
	}

	// 3. 执行爬取网站内容
	/* 不同网站的爬取函数不一样，需要根据实际情况实现
	// 目前先统一用 HuanqiuCrawler
	比如： 环球网 huanqiuwang 使用函数 CrawlHtml
	*/
	var newsMap map[string][]model.News
	for _, new := range allNewsData {
		logger.Debug("new:", new)
		logger.Debug("link,prefix:", new.Link, strings.Split(new.Link, ".")[1])
		newList := crawler.HuanQiuCrawer(new.Link, strings.Split(new.Link, ".")[1])
		newsMap[new.Source] = newList
		logger.Debug(newsMap)
	}

	// 5. 获取用户邮箱并推送到邮件队列
	user, err := getUserInfo(userId)
	if err != nil {
		logger.Errorf("获取用户邮箱失败: %v", err)
		return
	}

	// 3. 调用 Qwen API 进行内容分析
	analysisResult, err := analyzeNewsWithQwen(allNewsData, user)
	if err != nil {
		logger.Errorf("Qwen 分析失败: %v", err)
		return
	}

	// 4. 保存分析结果到数据库
	err = saveAnalysisResult(userId, taskId, analysisResult, allNewsData)
	if err != nil {
		logger.Errorf("保存分析结果失败: %v", err)
		return
	}

	// 创建邮件任务
	emailTask := EmailTask{
		UserID:   cast.ToUint64(userId),
		TaskID:   cast.ToUint64(taskId),
		Email:    user.Email,
		Subject:  fmt.Sprintf("新闻分析报告 - %s", time.Now().Format("2006-01-02")),
		Content:  analysisResult,
		NewsData: string(mustMarshal(allNewsData)),
	}

	// 推送到邮件队列
	err = EmailQueue.Produce(&emailTask)
	if err != nil {
		logger.Errorf("推送邮件任务失败: %v", err)
	} else {
		logger.Infof("邮件任务推送成功: UserID=%s, TaskID=%s", userId, taskId)
	}
}

// sendEmail 发送邮件
func sendEmail(emailTask EmailTask) error {
	logger.Infof("发送邮件: To=%s, Subject=%s", emailTask.Email, emailTask.Subject)

	// 格式化邮件内容
	htmlContent := smtp.FormatNewsAnalysisEmail(emailTask.Content, emailTask.NewsData)

	// 发送邮件
	err := smtp.SendEmail(emailTask.Email, emailTask.Subject, htmlContent)
	if err != nil {
		return fmt.Errorf("邮件发送失败: %v", err)
	}

	return nil
}

// mustMarshal JSON序列化辅助函数
func mustMarshal(v interface{}) []byte {
	data, err := json.Marshal(v)
	if err != nil {
		logger.Errorf("JSON marshal error: %v", err)
		return []byte("{}")
	}
	return data
}

func toCronExpression(timeStr string) (string, error) {
	// timeStr: "09:00"
	parts := strings.Split(timeStr, ":")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid time format")
	}
	hour, _ := strconv.Atoi(parts[0])
	minute, _ := strconv.Atoi(parts[1])

	// 生成 cron: "分钟 小时 * * *"
	return fmt.Sprintf("%d %d * * *", minute, hour), nil
}

// convertNewsSlice 转换新闻切片类型
func convertNewsSlice(newsData []*model.News) []model.News {
	result := make([]model.News, len(newsData))
	for i, news := range newsData {
		if news != nil {
			result[i] = *news
		}
	}
	return result
}
