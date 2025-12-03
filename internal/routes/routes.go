package routes

import (
	"go-sosmed-app/internal/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")

	// Auth
	api.POST("/register", controllers.Register)
	api.POST("/login", controllers.Login)

}
