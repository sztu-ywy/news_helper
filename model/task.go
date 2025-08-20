package model

type Task struct {
	Model
	TaskID   string `json:"task_id" gorm:"index"`
	TaskType string `json:"task_type" gorm:"index"`
	Status   string `json:"status"`
	UserID   uint64 `json:"user_id" gorm:"index"`
	User     *User  `json:"user" gorm:"foreignKey:UserID"`
}
