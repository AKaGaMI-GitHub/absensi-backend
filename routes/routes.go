package routes

import (
	"absen-backend/controller"

	"github.com/gin-gonic/gin"
)

func MainRoutes(router *gin.Engine) {
	router.POST("/login", controller.SignUp)

	router.Static("/static", "./public/uploads")
	api := router.Group("/api")
	{
		api.GET("/users", controller.GetUsers)
		api.POST("/users/create", controller.StoreUser)
	}
}
