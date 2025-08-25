package model

type Agent struct {
	Model
	Name   string `json:"name" gorm:"index" cosy:"add:required;update:omitempty"`
	Key    string `json:"key" cosy:"add:required;update:omitempty"`
	Remark string `json:"remark" cosy:"add:omitempty;update:omitempty"`
}
