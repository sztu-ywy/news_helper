package settings

import (
	"log"

	"gopkg.in/ini.v1"
)

var cfg *ini.File

func init() {
	Setup()
	Register("server", ServerSetting)
	Register("database", DataBaseSettings)
	Register("redis", RedisSettings)
	Register("georegeo", GeoregeoSettings)
	Register("jwt", JwtSettings)

}

var err error

func Setup() {
	cfg, err = ini.Load("app.ini")
	if err != nil {
		log.Fatalf("Failed to parse app.ini: %v", err)
	}

}
func Register(Name string, v interface{}) {
	mapTo(Name, v)
}

// 通用映射方法
func mapTo(section string, v interface{}) {
	err := cfg.Section(section).MapTo(v)
	if err != nil {
		log.Fatalf("Cfg.MapTo %s err: %v", section, err)
	}
}
