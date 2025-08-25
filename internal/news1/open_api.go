package news1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"news_helper/model"
	newsSettings "news_helper/settings"
	"strings"
	"time"

	"github.com/uozi-tech/cosy/logger"
)

// QwenRequest Qwen API 请求结构
type QwenRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

// Message 消息结构
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// QwenResponse Qwen API 响应结构
type QwenResponse struct {
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

// Choice 选择结构
type Choice struct {
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

// Usage 使用量统计
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// analyzeNewsWithQwen 使用 Qwen API 分析新闻
func analyzeNewsWithQwen(newsData []*model.News, user *model.User) (string, error) {
	logger.Info("开始使用 Qwen 分析新闻内容")

	// 检查 API Key 配置
	if newsSettings.QwenSettings.ApiKey == "" {
		logger.Warn("Qwen API Key 未配置，使用模拟数据")
		return generateMockAnalysis(newsData), nil
	}

	// 构建完善的分析提示词
	prompt := buildAnalysisPrompt(newsData, user)

	// 调用 Qwen API
	analysisResult, err := callQwenAPI(prompt)
	if err != nil {
		logger.Errorf("调用 Qwen API 失败: %v", err)
		// 如果 API 调用失败，返回模拟数据
		return generateMockAnalysis(newsData), nil
	}

	return analysisResult, nil
}

// buildAnalysisPrompt 构建分析提示词
func buildAnalysisPrompt(newsData []*model.News, user *model.User) string {
	var prompt strings.Builder

	// 系统角色设定
	prompt.WriteString("你是一位专业的金融分析师和投资顾问，具有丰富的市场分析经验。请基于以下新闻内容，为用户提供专业的投资分析报告。\n\n")

	// 用户信息（如果有的话）
	if user != nil {
		prompt.WriteString(fmt.Sprintf("用户信息：%s\n\n", user.Name))
	}

	// 分析要求
	prompt.WriteString("分析要求：\n")
	prompt.WriteString("1. 市场趋势分析：基于新闻内容分析当前市场趋势和方向\n")
	prompt.WriteString("2. 投资机会识别：识别潜在的投资机会和热点板块\n")
	prompt.WriteString("3. 风险评估：分析可能的市场风险和注意事项\n")
	prompt.WriteString("4. 具体建议：提供具体的投资建议和操作策略\n")
	prompt.WriteString("5. 时间维度：区分短期、中期、长期的投资视角\n\n")

	// 新闻内容
	prompt.WriteString("新闻内容：\n")
	for i, news := range newsData {
		prompt.WriteString(fmt.Sprintf("【新闻 %d】\n", i+1))
		prompt.WriteString(fmt.Sprintf("标题：%s\n", news.Title))
		prompt.WriteString(fmt.Sprintf("来源：%s\n", news.Source))
		prompt.WriteString(fmt.Sprintf("发布时间：%s\n", time.Unix(int64(news.Time), 0).Format("2006-01-02 15:04:05")))
		prompt.WriteString(fmt.Sprintf("内容：%s\n\n", news.Content))
	}

	// 输出格式要求
	prompt.WriteString("请按以下格式输出分析报告：\n")
	prompt.WriteString("# 新闻分析报告\n\n")
	prompt.WriteString("## 📊 市场趋势分析\n")
	prompt.WriteString("[在此分析市场整体趋势]\n\n")
	prompt.WriteString("## 🎯 投资机会\n")
	prompt.WriteString("[在此识别投资机会和热点板块]\n\n")
	prompt.WriteString("## ⚠️ 风险提示\n")
	prompt.WriteString("[在此分析潜在风险]\n\n")
	prompt.WriteString("## 💡 投资建议\n")
	prompt.WriteString("### 短期建议（1-3个月）\n")
	prompt.WriteString("[短期投资建议]\n\n")
	prompt.WriteString("### 中期建议（3-12个月）\n")
	prompt.WriteString("[中期投资建议]\n\n")
	prompt.WriteString("### 长期建议（1年以上）\n")
	prompt.WriteString("[长期投资建议]\n\n")
	prompt.WriteString("## 📈 关键指标关注\n")
	prompt.WriteString("[建议关注的关键指标和数据]\n\n")

	return prompt.String()
}

// callQwenAPI 调用 Qwen API
func callQwenAPI(prompt string) (string, error) {
	// Qwen API 端点
	apiURL := "https://dashscope.aliyuncs.com/api/v1/services/aigc/text-generation/generation"

	// 构建请求数据
	requestData := QwenRequest{
		Model: "qwen-turbo", // 使用 qwen-turbo 模型
		Messages: []Message{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	// 序列化请求数据
	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return "", fmt.Errorf("序列化请求数据失败: %v", err)
	}

	// 创建 HTTP 请求
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+newsSettings.QwenSettings.ApiKey)

	// 发送请求
	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %v", err)
	}

	// 检查 HTTP 状态码
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API 请求失败，状态码: %d, 响应: %s", resp.StatusCode, string(body))
	}

	// 解析响应
	var qwenResp QwenResponse
	if err := json.Unmarshal(body, &qwenResp); err != nil {
		return "", fmt.Errorf("解析响应失败: %v", err)
	}

	// 检查响应是否有效
	if len(qwenResp.Choices) == 0 {
		return "", fmt.Errorf("API 响应无效，没有返回内容")
	}

	// 记录使用量
	logger.Infof("Qwen API 调用成功，使用 tokens: %d", qwenResp.Usage.TotalTokens)

	return qwenResp.Choices[0].Message.Content, nil
}

// generateMockAnalysis 生成模拟分析结果
func generateMockAnalysis(newsData []*model.News) string {
	var analysis strings.Builder

	analysis.WriteString("# 新闻分析报告\n\n")
	analysis.WriteString(fmt.Sprintf("*基于 %d 条新闻的分析报告（模拟数据）*\n\n", len(newsData)))

	analysis.WriteString("## 📊 市场趋势分析\n")
	analysis.WriteString("根据当前新闻内容分析，市场整体呈现积极态势。主要表现在：\n")
	analysis.WriteString("- 政策环境相对稳定\n")
	analysis.WriteString("- 行业发展动力充足\n")
	analysis.WriteString("- 市场情绪趋于乐观\n\n")

	analysis.WriteString("## 🎯 投资机会\n")
	analysis.WriteString("基于新闻分析，以下领域值得关注：\n")
	analysis.WriteString("- 科技创新板块：人工智能、新能源等\n")
	analysis.WriteString("- 消费升级相关：高端制造、品牌消费\n")
	analysis.WriteString("- 基础设施建设：新基建、绿色能源\n\n")

	analysis.WriteString("## ⚠️ 风险提示\n")
	analysis.WriteString("投资过程中需要注意以下风险：\n")
	analysis.WriteString("- 市场波动风险：短期内可能存在调整\n")
	analysis.WriteString("- 政策变化风险：关注相关政策动向\n")
	analysis.WriteString("- 行业竞争风险：注意行业集中度变化\n\n")

	analysis.WriteString("## 💡 投资建议\n")
	analysis.WriteString("### 短期建议（1-3个月）\n")
	analysis.WriteString("- 保持谨慎乐观态度\n")
	analysis.WriteString("- 关注热点板块轮动机会\n")
	analysis.WriteString("- 控制仓位，分散投资\n\n")

	analysis.WriteString("### 中期建议（3-12个月）\n")
	analysis.WriteString("- 重点关注成长性行业\n")
	analysis.WriteString("- 选择基本面良好的优质标的\n")
	analysis.WriteString("- 适当增加配置比例\n\n")

	analysis.WriteString("### 长期建议（1年以上）\n")
	analysis.WriteString("- 坚持价值投资理念\n")
	analysis.WriteString("- 关注行业龙头企业\n")
	analysis.WriteString("- 定期调整投资组合\n\n")

	analysis.WriteString("## 📈 关键指标关注\n")
	analysis.WriteString("建议重点关注以下指标：\n")
	analysis.WriteString("- 宏观经济指标：GDP、CPI、PMI\n")
	analysis.WriteString("- 行业指标：相关行业景气度指数\n")
	analysis.WriteString("- 技术指标：市场成交量、资金流向\n\n")

	analysis.WriteString("---\n")
	analysis.WriteString("*本报告仅供参考，投资有风险，决策需谨慎*")

	return analysis.String()
}
