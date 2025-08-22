package news1

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"news_helper/model"
	"runtime"
	"time"

	"github.com/spf13/cast"
	"github.com/uozi-tech/cosy/queue"
	"github.com/uozi-tech/cosy/redis"

	"github.com/uozi-tech/cosy/logger"
)

const (
	CrawlerRedisPrefix = "crawler"
	CrawlerExpiredTime = 24 * 30 * time.Hour
)

var (
	TaskQueue        *queue.Queue[model.Task]
	NewsCrawlerQueue *queue.Queue[model.News]
)

func InitQueue(ctx context.Context) {
	logger.Debug("222222222")
	NewsCrawlerQueue = queue.New[model.News]("news_crawler", queue.LeftToRight)
	TaskQueue = queue.New[model.Task]("task_queue", queue.LeftToRight)

	go ConsumeTaskQueue(ctx)

}

func CrawlHtml(method, url string, headers http.Header, queryParams url.Values, data interface{}) ([]byte, error) {
	var reqBody io.Reader
	if data != nil {
		jsonData, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewReader(jsonData)
	}
	// 1. 发起 HTTP 请求
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, err
	}
	req.Header = headers
	req.URL.RawQuery = queryParams.Encode()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Errorf("请求失败:", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Errorf("HTTP 错误: %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
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
				a, err = redis.HSet(str_userId, str_taskId, jsonData)
				if err != nil {
					logger.Errorf("2HSet error: %v,a: %v", err, a)
				}
				err = redis.Set(str_taskId, task.DailyTime, CrawlerExpiredTime)
				if err != nil {
					logger.Errorf("3HSet error: %v,a: %v", err, a)
				}
			}()
		}
	}
}
