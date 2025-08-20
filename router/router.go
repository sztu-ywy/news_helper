package router

import (
	"github.com/gin-gonic/gin"
)

func Init() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Hello, World!"})
	})
	r.Run("0.0.0.0:8084")
}

// var r *gin.Engine

// // Init initializes the router, expose this function to be able to reset the testing router
// func Init() {
// 	r = gin.New()

// 	r.Use(gin.Recovery())

// 	r.Use(gin.Logger())

// 	r.NoRoute(func(c *gin.Context) {
// 		c.JSON(http.StatusNotFound, gin.H{
// 			"message": "not found",
// 		})
// 	})
// }

// func GetEngine() *gin.Engine {
// 	return r
// }
