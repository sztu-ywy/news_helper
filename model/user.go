package model

import (
	"encoding/json"
	"time"

	"news_helper/internal/acl"
	"news_helper/internal/helper"

	"github.com/uozi-tech/cosy/logger"
	"github.com/uozi-tech/cosy/redis"
	"gorm.io/gorm"
)

const (
	UserStatusActive = 1
	UserStatusBan    = -1
)

type User struct {
	Model
	Name        string     `json:"name,omitempty" cosy:"add:required;update:omitempty;list:fussy"`
	Password    string     `json:"-" cosy:"json:password;add:required;update:omitempty;delete:omitempty;get:omitempty;"`
	Email       string     `json:"email,omitempty" cosy:"add:required;update:omitempty;list:fussy" gorm:"type:varchar(255);index"`
	Phone       string     `json:"phone,omitempty" cosy:"all:omitempty;list:fussy" gorm:"index"`
	LastActive  int64      `json:"last_active,omitempty"`
	UserGroupID uint64     `json:"user_group_id,omitempty" cosy:"all:omitempty;list:eq" gorm:"index;default:0"`
	UserGroup   *UserGroup `json:"user_group,omitempty" cosy:"item:preload;list:preload"`
	//1是启用，-1是禁用
	Status int `json:"status,omitempty" cosy:"add:min=-1,max=1;update:omitempty,min=-1,max=1;list:in" gorm:"default:1"`
}

func (u *User) AfterUpdate(_ *gorm.DB) (err error) {
	if u.ID == 0 {
		logger.Warn("the after update hook of user model detected an invalid user id(0), " +
			"this will not clean the cache of the user you expected")
		return
	}
	key := helper.BuildUserKey(u.ID)
	err = redis.Del(key)
	if err != nil {
		logger.Error(err)
		return nil
	}
	return
}

// UpdateLastActive update last active time in redis
func (u *User) UpdateLastActive() (now int64) {
	now = time.Now().Unix()
	key := helper.BuildUserKey(u.ID, "last_active")
	_ = redis.Set(key, now, 0)
	return
}

// GetUserGroup get user group from cache or db
func (u *User) GetUserGroup() (group *UserGroup, err error) {
	if u.UserGroupID == 0 {
		return &UserGroup{}, nil
	}
	// Check cache
	key := getUserGroupKey(u.UserGroupID)
	value, err := redis.Get(key)
	group = &UserGroup{}
	if err != nil || value == "" {
		// query group and build permissions map and set to cache
		err := db.First(&group, u.UserGroupID).Error

		if err != nil {
			return nil, err
		}

		bytes, err := json.Marshal(group)
		if err != nil {
			return nil, err
		}

		err = redis.Set(key, string(bytes), 5*time.Minute)
		if err != nil {
			return nil, err
		}
		return group, nil
	}

	bytes := []byte(value)
	err = json.Unmarshal(bytes, group)

	if err != nil {
		return nil, err
	}

	return
}

// GetPermissionsMap get permissions map from user group
func (u *User) GetPermissionsMap() (permissionsMap acl.Map) {
	group, err := u.GetUserGroup()
	if err != nil {
		return
	}
	permissionsMap = group.Permissions.ToMap()
	return
}

// IsAdmin check whether user is admin
func (u *User) IsAdmin() bool {
	return acl.Can(u.GetPermissionsMap(), acl.All, acl.Write)
}

// Can check whether the user can do the action
func (u *User) Can(subject acl.Subject, action acl.Action) bool {
	return acl.Can(u.GetPermissionsMap(), subject, action)
}
