package model

import (
	"encoding/json"
	"strings"

	"github.com/spf13/cast"
	"github.com/uozi-tech/cosy/redis"
)

type Setting struct {
	Model
	Name string `json:"name"`
	Meta string `json:"meta"`
}

func (s *Setting) Marshal(meta interface{}) {
	bytes, _ := json.Marshal(meta)

	s.Meta = string(bytes)

	return
}

func (s *Setting) Unmarshal() (meta interface{}) {
	_ = json.Unmarshal([]byte(s.Meta), &meta)
	return
}

func (s *Setting) UnmarshalTo(dest interface{}) {
	_ = json.Unmarshal([]byte(s.Meta), dest)
	return
}

func (s *Setting) Insert() {
	db.Create(s)
}

func (s *Setting) Save() error {
	return db.Save(s).Error
}

func InitRuntimeSettings() {
	var settings []Setting
	db.Find(&settings)

	for _, v := range settings {
		key := buildSettingKey(v.Name)
		_ = redis.Set(key, v.Meta, 0)
	}
}

func buildSettingKey(key interface{}) string {
	var sb strings.Builder
	sb.WriteString("settings:")
	sb.WriteString(cast.ToString(key))
	return sb.String()
}
