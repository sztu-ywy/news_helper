package model

import (
	"strings"

	"news_helper/internal/acl"
	"github.com/spf13/cast"
	"github.com/uozi-tech/cosy/redis"
)

type UserGroup struct {
	Model
	Name        string             `json:"name" cosy:"add:required;update:omitempty;list:fussy" gorm:"type:varchar(255);uniqueIndex"`
	Permissions acl.PermissionList `json:"permissions" cosy:"all:required" gorm:"type:longtext;serializer:json"`
}

func getUserGroupKey(id uint64) string {
	var sb strings.Builder
	sb.WriteString("user_group:")
	sb.WriteString(cast.ToString(id))
	return sb.String()
}

func (u *UserGroup) CleanCache() {
	_ = redis.Del(getUserGroupKey(u.ID))
}
