package model

type Task struct {
	Model
	UserID    uint64   `json:"user_id" cosy:"add:required;update:omitempty" gorm:"index"`
	User      *User    `json:"user" cosy:"item:preload;list:preload" gorm:"foreignKey:UserID"`
	MediaIDs  []uint64 `json:"media_ids" cosy:"add:omitempty;update:omitempty" gorm:"serializer:json"` // 修复JSON标签，允许绑定media_ids
	Medias    []*Media `json:"medias" cosy:"item:preload;list:preload" gorm:"many2many:task_medias;"`
	DailyTime string   `json:"daily_time" cosy:"add:omitempty;update:omitempty"`
	Remark    string   `json:"remark" cosy:"add:omitempty;update:omitempty"`
	TaskType  string   `json:"task_type" cosy:"add:omitempty;update:omitempty" gorm:"index"`
	Status    string   `json:"status" cosy:"add:omitempty;update:omitempty"`
}
