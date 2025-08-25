package main

import (
	"flag"

	"news_helper/internal/audit"
	"news_helper/internal/limiter"
	"news_helper/internal/news1"
	"news_helper/model"

	"news_helper/model/view"
	"news_helper/query"
	"news_helper/router"

	"github.com/uozi-tech/cosy"
	mysql "github.com/uozi-tech/cosy-driver-mysql"
	"github.com/uozi-tech/cosy/settings"
)

type Config struct {
	ConfPath string
	Maintain string
}

var cfg Config

func init() {
	// 指定默认配置文件，命令行与参数
	flag.StringVar(&cfg.ConfPath, "config", "app.ini", "Specify the configuration file")
	// 解析
	flag.Parse()
}

func main() {
	// 注册模型，用于数据库迁移和查询等
	cosy.RegisterModels(model.GenerateAllModel()...)
	// 注册初始化函数，数据库连接实例化、创建admin用户，email:admin;password:admin初始化查询模块、实例化数据库对象传递给模型。创建数据库视图，初始化限流器，控制请求速率
	cosy.RegisterInitFunc(func() {
		db := cosy.InitDB(mysql.Open(settings.DataBaseSettings))
		query.Init(db)
		model.Use(db)
		view.CreateViews(db)
		limiter.Init()
	},
		// 初始化运行时配置，初始化路由
		// model.InitRuntimeSettings,
		router.InitRouter,
	)
	// 注册后台运行函数，用于异步运行收集日志和监控
	cosy.RegisterGoroutine(
		audit.Init,
		news1.InitNews,
		news1.InitQueue, // InitQueue 现在接受 context 参数，会在 cosy 框架中自动传递

	)

	// cosy启动
	cosy.Boot(cfg.ConfPath)
	// tmpemail := "123@qq.com"
	// Register1(tmpemail)

}
