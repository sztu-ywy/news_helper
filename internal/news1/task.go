package news1

import (
	"context"
	"encoding/json"
	"fmt"
	"news_helper/internal/queue"
	"news_helper/model"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/robfig/cron"
	"github.com/spf13/cast"
	"github.com/uozi-tech/cosy/logger"
	"github.com/uozi-tech/cosy/redis"
)

const (
	CrawlerRedisPrefix = "crawler"
	CrawlerExpiredTime = 24 * 30 * time.Hour
	CrawlerUser        = "crawler:user:"
	CrawlerTask        = "crawler:task:"
)

var (
	TaskQueue        *queue.Queue[model.Task]
	NewsCrawlerQueue *queue.Queue[model.News]
)

func InitQueue(ctx context.Context) {
	NewsCrawlerQueue = queue.New[model.News]("news_crawler", queue.LeftToRight)
	TaskQueue = queue.New[model.Task]("task_queue", queue.LeftToRight)

	go ConsumeTaskQueue(ctx)

}
func ConsumeTaskQueue(ctx context.Context) {

	ch, err := TaskQueue.Subscribe(ctx)
	if err != nil {
		logger.Errorf("consume task queue error: %v", err)
		return
	}
	for {
		select {
		case <-ctx.Done():
			return
		case task := <-ch:
			logger.Infof("task: %v", task)
			// todo
			func() {
				defer func() {
					if err := recover(); err != nil {
						buf := make([]byte, 1024)
						runtime.Stack(buf, false)
						logger.Errorf("consume task queue error: %v,buf: %s", err, buf)
					}
				}()
				str_taskId := cast.ToString(task.ID)
				str_userId := cast.ToString(task.UserID)
				a, err := redis.HSet(CrawlerRedisPrefix, str_taskId, str_userId)
				if err != nil {
					logger.Errorf("1HSet error: %v,a: %v", err, a)
				}
				jsonData, err := json.Marshal(task.MediaIDs)
				if err != nil {
					logger.Errorf("json marshal error: %v", err)
				}
				a, err = redis.HSet(CrawlerUser+str_userId, str_taskId, jsonData)
				if err != nil {
					logger.Errorf("2HSet error: %v,a: %v", err, a)
				}
				err = redis.Set(CrawlerTask+str_taskId, task.DailyTime, CrawlerExpiredTime)
				if err != nil {
					logger.Errorf("3HSet error: %v,a: %v", err, a)
				}
			}()
		}
	}
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
	logger.Info("doTask")
	logger.Info("userId, taskIds:", userId, taskIds)

	for task, time := range tasks_time {

		RegisterTask(ctx, c, userId, task, time)
	}

}

func RegisterTask(ctx context.Context, c *cron.Cron, userId, taskId string, taskTime string) {
	// todo
	// _, err := cron.AddFunc("0 1,9,11 * * *", c.Clean)
	cronTime, _ := toCronExpression(taskTime)
	c.AddFunc(cronTime, func() {
		user_id := userId
		task_id := taskId
		doTask(user_id, task_id)
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

func toCronExpression(timeStr string) (string, error) {
	// timeStr: "09:00"
	parts := strings.Split(timeStr, ":")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid time format")
	}
	hour, _ := strconv.Atoi(parts[0])
	minute, _ := strconv.Atoi(parts[1])

	// 生成 cron: "分钟 小时 * * *"
	return fmt.Sprintf("%d %d * * *", minute, hour), nil
}
