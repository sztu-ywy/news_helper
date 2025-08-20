package main

import (
	"flag"
	"news_helper/router"
	"news_helper/settings"
)

type Config struct {
	ConfPath string
	Maintain string
}

var cfg Config

func init() {
	flag.StringVar(&cfg.ConfPath, "config", "app.ini", "Specify the configuration file")
	flag.Parse()
}

func main() {
	settings.InitDatabase(*settings.DataBaseSettings)
	settings.InitRedis(*settings.RedisSettings)
	router.Init()

}
