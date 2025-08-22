package news1

import (
	"context"
	"encoding/json"
	"news_helper/internal/queue"
	"news_helper/model"
	"runtime"
	"time"

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
