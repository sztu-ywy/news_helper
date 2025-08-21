package model

type News struct {
	Model
	MediaId uint   `json:"media_id" gorm:"index"`
	Title   string `json:"title" gorm:"index"`
	Content string `json:"content"`
	Link    string `json:"link"`
	Source  string `json:"source"`
	Time    uint64 `json:"time"`
}
