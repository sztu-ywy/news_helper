package router

import (
	"news_helper/api/admin/task"

	"git.uozi.org/uozi/burn-api/api/global"
	"github.com/gin-gonic/gin"
	"github.com/uozi-tech/cosy"
)

func InitRouter() {
	r := cosy.GetEngine()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Hello, World!"})
	})
	global.InitRouter(r)

	taskRouter := r.Group("admin")
	{
		task.InitRouter(taskRouter)
	}

}
