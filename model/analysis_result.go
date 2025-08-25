package model

// AnalysisResult 新闻分析结果模型
type AnalysisResult struct {
	Model
	UserID       uint64 `json:"user_id" cosy:"add:required;update:omitempty" gorm:"index"`
	User         *User  `json:"user" cosy:"item:preload;list:preload" gorm:"foreignKey:UserID"`
	TaskID       uint64 `json:"task_id" cosy:"add:required;update:omitempty" gorm:"index"`
	Task         *Task  `json:"task" cosy:"item:preload;list:preload" gorm:"foreignKey:TaskID"`
	Title        string `json:"title" cosy:"add:required;update:omitempty"`                    // 分析报告标题
	Content      string `json:"content" cosy:"add:required;update:omitempty" gorm:"type:text"` // 分析内容
	NewsCount    int    `json:"news_count" cosy:"add:omitempty;update:omitempty"`              // 分析的新闻数量
	AnalysisDate uint64 `json:"analysis_date" cosy:"add:omitempty;update:omitempty"`          // 分析日期
	Status       string `json:"status" cosy:"add:omitempty;update:omitempty" gorm:"default:pending"` // 状态：pending, completed, failed
}
