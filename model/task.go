package model

type Task struct {
	Model
	TaskType string   `json:"task_type" gorm:"index"`
	Status   string   `json:"status"`
	UserID   uint64   `json:"user_id" gorm:"index"`
	User     *User    `json:"user" gorm:"foreignKey:UserID"`
	MediaIDs []uint64 `json:"-" gorm:"serializer:json"`
	Medias   []*Media `json:"medias" gorm:"many2many:task_medias;"`
}
