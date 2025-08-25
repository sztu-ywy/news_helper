package smtp

import (
	"fmt"
	"net/smtp"
	"strings"

	"github.com/uozi-tech/cosy/logger"
	"news_helper/settings"
)

// EmailConfig 邮件配置
type EmailConfig struct {
	Host     string
	Port     string
	Email    string
	Password string
}

// SendEmail 发送邮件
func SendEmail(to, subject, body string) error {
	// 获取邮件配置
	config := EmailConfig{
		Host:     settings.EmailSettings.Host,
		Port:     settings.EmailSettings.Port,
		Email:    settings.EmailSettings.Email,
		Password: settings.EmailSettings.Password,
	}

	// 检查配置是否完整
	if config.Host == "" || config.Email == "" || config.Password == "" {
		logger.Warn("邮件配置不完整，跳过发送")
		return fmt.Errorf("邮件配置不完整")
	}

	// 构建邮件内容
	msg := []string{
		fmt.Sprintf("From: %s", config.Email),
		fmt.Sprintf("To: %s", to),
		fmt.Sprintf("Subject: %s", subject),
		"MIME-Version: 1.0",
		"Content-Type: text/html; charset=UTF-8",
		"",
		body,
	}

	// SMTP 认证
	auth := smtp.PlainAuth("", config.Email, config.Password, config.Host)

	// 发送邮件
	addr := fmt.Sprintf("%s:%s", config.Host, config.Port)
	err := smtp.SendMail(addr, auth, config.Email, []string{to}, []byte(strings.Join(msg, "\r\n")))
	if err != nil {
		return fmt.Errorf("发送邮件失败: %v", err)
	}

	logger.Infof("邮件发送成功: %s -> %s", config.Email, to)
	return nil
}

// FormatNewsAnalysisEmail 格式化新闻分析邮件内容
func FormatNewsAnalysisEmail(analysisResult, newsData string) string {
	html := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>新闻分析报告</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; margin: 20px; }
        .header { background-color: #f4f4f4; padding: 20px; border-radius: 5px; }
        .content { margin: 20px 0; }
        .analysis { background-color: #e8f4fd; padding: 15px; border-radius: 5px; margin: 10px 0; }
        .news-data { background-color: #f9f9f9; padding: 15px; border-radius: 5px; margin: 10px 0; }
        .footer { color: #666; font-size: 12px; margin-top: 30px; }
    </style>
</head>
<body>
    <div class="header">
        <h1>📊 新闻分析报告</h1>
        <p>生成时间: %s</p>
    </div>
    
    <div class="content">
        <h2>🔍 分析结果</h2>
        <div class="analysis">
            <pre>%s</pre>
        </div>
        
        <h2>📰 相关新闻数据</h2>
        <div class="news-data">
            <details>
                <summary>点击查看详细新闻数据</summary>
                <pre>%s</pre>
            </details>
        </div>
    </div>
    
    <div class="footer">
        <p>此邮件由新闻助手系统自动生成</p>
    </div>
</body>
</html>
`,
		fmt.Sprintf("%s", "2024-01-01 12:00:00"), // 这里应该传入实际时间
		analysisResult,
		// newsData,
	)

	return html
}
