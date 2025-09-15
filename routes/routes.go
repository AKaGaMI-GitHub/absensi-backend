package routes

import (
	"github.com/gin-gonic/gin"
)

func MainRoutes(router *gin.Engine) {
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "hello world",
		})
	})
}
