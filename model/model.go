package model

import (
	"news_helper/internal/helper"

	"github.com/spf13/cast"
	"github.com/uozi-tech/cosy/logger"
	"github.com/uozi-tech/cosy/redis"
	"gorm.io/gen"
	"gorm.io/gorm"
	"gorm.io/plugin/soft_delete"
)

var db *gorm.DB

type Model struct {
	ID        uint64                `gorm:"primary_key" json:"id,string"`
	CreatedAt uint64                `json:"created_at,omitempty" gorm:"autoCreateTime" cosy:"list:between"`
	UpdatedAt uint64                `json:"updated_at,omitempty" gorm:"autoUpdateTime" cosy:"list:between"`
	DeletedAt soft_delete.DeletedAt `json:"deleted_at,omitempty" gorm:"index;default:0"`
}

// main.go执行了GenerateAllModel进行自动迁移，还要在GenerateAllModel()中注册
func GenerateAllModel() []any {
	return []any{
		User{},
		UserGroup{},
		Upload{},
		Task{},  // 添加 Task 模型注册
		Media{}, // 添加 Media 模型注册
		News{},  // 添加 News 模型注册
	}
}

func Use(tx *gorm.DB) {
	db = tx
}

type Method interface {
	// FirstByID Where("id=@id")
	FirstByID(id uint64) (*gen.T, error)
	// DeleteByID update @@table set deleted_at=NOW() where id=@id
	DeleteByID(id uint64) error
}

func DropCache(prefix string, ids ...uint64) {
	keys := make([]string, 0)
	for _, id := range ids {
		if id != 0 {
			keys = append(keys, helper.BuildKey(prefix, cast.ToString(id)))
		}
	}
	err := redis.Del(keys...)
	if err != nil {
		logger.Debug(err)
	}
}
