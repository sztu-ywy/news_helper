package settings

import (
	_ "github.com/spf13/viper/remote" // 导入 INI 驱动
)

type DataBase struct {
	Type        string `json:"type"`
	Host        string `json:"host"`
	Port        uint   `json:"port"`
	User        string `json:"user"`
	Password    string `json:"-,omitempty"`
	Name        string `json:"name"`
	TablePrefix string `json:"table_prefix"`
}

var DataBaseSettings = &DataBase{}
