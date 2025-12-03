package controllers

import (
	"go-sosmed-app/internal/utils"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(c *gin.Context) {
	userId, err := utils.ValidateJWT(c)

	if err != nil {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	c.Set("userId", userId)
	c.Next()
}
