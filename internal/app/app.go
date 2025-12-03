package app

import (
	"go-sosmed-app/internal/config"
	"go-sosmed-app/internal/db"
	"go-sosmed-app/internal/routes"

	"github.com/gin-gonic/gin"
)

func Run() {
	config.LoadEnv()

	db.Init()

	router := gin.Default()
	routes.RegisterRoutes(router)
	router.Run(":" + config.Getenv("API_PORT", "8080"))
}
