package news1

import (
	"context"
	"sync"
	"time"

	"github.com/robfig/cron"

	"github.com/uozi-tech/cosy/logger"
	"github.com/uozi-tech/cosy/redis"
)

func InitNews(ctx context.Context) {

	redis.HSet("crawler", "123", "0312")
	redis.HSet("crawler", "1234", "1222")
	time.Sleep(1 * time.Second)
	logger.Info("InitNews")

	hashMap, err := redis.HGetAll("crawler")
	if err != nil {
		logger.Error("HGetAll", err)
		return
	}
	userTaskMap := make(map[string][]string)
	for k, v := range hashMap {
		userTaskMap[v] = append(userTaskMap[v], k)
	}

	var wg sync.WaitGroup
	for k, v := range userTaskMap {
		logger.Info(k, v)
		wg.Add(1)
		go func(k string, v []string) {
			defer wg.Done()
			// todo
			TaskTimer(ctx, k, v)
		}(k, v)

	}
	//  查询 redis中 key为crawler 的哈希表
	// 表中key为任务 id，value为用户的 id
	// 每个用户开启一个协程
	// 把所有的不同用户的任务 id 放到不同协程的数组 中

	// 查询 任务 id 为 key，value 为任务的时间戳，

	// key 为 userId，key_为任务 id，value 为媒体集合

	// 如果集合为空，则返回
	//如果集合不为空，则遍历集合，缓存，依次爬取列表中新闻名称对应的链接，内容储存到数据库中，并整合内容，调用 api 接口，总结投资意向，存到数据库中

	//搜索数据库，通过邮箱，推送到指定的用户

}
func TaskTimer(ctx context.Context, userId string, taskIds []string) {
	tasks_time, err := redis.HGetAll(userId)
	if err != nil {
		logger.Error("HGetAll", err)
		return
	}
	// todo
	// tasksTime := make([]string, 0)
	c := cron.New()
	for task, time := range tasks_time {
		RegisterTask(ctx, c, userId, task, time)
	}

}

func RegisterTask(ctx context.Context, c *cron.Cron, userId, taskId string, taskTime string) {
	// todo
	// _, err := cron.AddFunc("0 1,9,11 * * *", c.Clean)
	c.AddFunc(taskTime, func() {
		logger.Info("doTask")
		doTask(userId, taskId)
	})

}

func doTask(userId, taskId string) {
	// todo
	newsList, err := redis.Get(taskId)
	if err != nil {
		logger.Error("Get", err)
		return
	}
	for _, news := range newsList {
		logger.Info(news)
		// todo
		// 爬取新闻内容
		// 调用 api 接口，总结投资意向
		// 存到数据库中
	}

	logger.Info(newsList)
}
