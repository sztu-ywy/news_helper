package task

import (
	"bytes"
	"io"
	"news_helper/internal/news1"

	"news_helper/model"

	"github.com/gin-gonic/gin"
	"github.com/uozi-tech/cosy"
	"github.com/uozi-tech/cosy/logger"
)

// redis中 key为crawler 的哈希表
// 表中key为任务 id，value为用户的 id

// key 为 userId，key_为任务 id，value 为媒体集合

// 查询 任务 id 为 key，value 为任务的时间戳，

func InitRouter(r *gin.RouterGroup) {
	logger.Debug("1111")
	c := cosy.Api[model.Task]("tasks")
	c.CreateHook(func(c *cosy.Ctx[model.Task]) {
		c.BeforeExecuteHook(func(c *cosy.Ctx[model.Task]) {
			logger.Debug("BeforeExecuteHook 开始执行")

			// 尝试直接从 gin.Context 获取原始请求数据
			ginCtx := c.Context
			// 读取原始请求体
			body, err := ginCtx.GetRawData()
			if err == nil {
				logger.Debugf("原始请求体: %s", string(body))

				// 重新设置请求体，因为 GetRawData 会消耗它
				ginCtx.Request.Body = io.NopCloser(bytes.NewBuffer(body))
			}

			logger.Debugf("请求数据绑定前 - OriginModel: %+v", c.OriginModel)
			logger.Debugf("请求数据绑定前 - Model: %+v", c.Model)
		})
		c.ExecutedHook(func(c *cosy.Ctx[model.Task]) {
			logger.Debug("ExecutedHook 开始执行")

			// 检查 TaskQueue 是否为 nil
			if news1.TaskQueue == nil {
				logger.Error("TaskQueue 为 nil")
				return
			}
			logger.Debug("TaskQueue 不为 nil")

			task := c.Model
			logger.Debugf("获取到 task: %+v", task)

			// 检查 task 是否为空
			if task.ID == 0 {
				logger.Error("task ID 为 0，可能未正确生成")
				return
			}

			// 必须确认 taskID 已经生成
			logger.Debug("准备调用 TaskQueue.Produce")
			err := news1.TaskQueue.Produce(&task)
			if err != nil {
				logger.Errorf("任务放入队列失败: %v", err)
			} else {
				logger.Debug("任务成功放入队列")
			}
		})

	})
	c.InitRouter(r)
}
