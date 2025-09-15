package main

import (
	"absen-backend/config"
	"absen-backend/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectDB()
	router := gin.Default()
	routes.MainRoutes(router)
	router.Run("localhost:8000")
}
