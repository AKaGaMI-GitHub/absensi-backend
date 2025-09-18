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
		//user
		api.GET("/users", controller.GetUsers)
		api.GET("/users/:username", controller.GetUserByUsername)
		api.POST("/users/create", controller.StoreUser)
		api.PATCH("/users/:uuid", controller.UpdateUser)
		api.DELETE("/users/:uuid", controller.DeleteUser)

		//role
		api.GET("/roles", controller.GetRoleUser)
		api.POST("/roles/create", controller.StoreRole)
		api.PATCH("/roles/:uuid", controller.UpdateRole)
		api.DELETE("/roles/:uuid", controller.DeleteRole)
	}
}
