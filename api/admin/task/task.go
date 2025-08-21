package task

import (
	"news_helper/model"

	"github.com/gin-gonic/gin"
	"github.com/uozi-tech/cosy"
)

func InitRouter(r *gin.RouterGroup) {
	c := cosy.Api[model.Task]("tasks")
	c.InitRouter(r)
}
