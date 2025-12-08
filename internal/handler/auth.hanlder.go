package handler

import (
	"errors"
	"go-sosmed-app/internal/db"
	"go-sosmed-app/internal/models"
	"go-sosmed-app/internal/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Register(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": utils.ParseErrorMessage(err)})
		return
	}

	var user models.User
	err := db.DB.Where("username = ?", req.Username).First(&user).Error

	if err == nil {
		c.JSON(409, gin.H{"error": "Username already taken"})
		return
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(500, gin.H{"error": "Failed to check if user exists"})
		return
	}

	user = models.User{Username: req.Username, Password: utils.HashPassword(req.Password)}
	if err := db.DB.Create(&user).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(201, gin.H{"message": "User created successfully"})
}

func Login(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": utils.ParseErrorMessage(err)})
		return
	}

	var user models.User

	if err := db.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(404, gin.H{"error": "User not found"})
			return
		}

		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if !utils.CheckPass(req.Password, user.Password) {
		c.JSON(401, gin.H{"error": "Password is incorrect"})
		return
	}

	token, err := utils.GenerateJWT(user.ID)

	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(200, gin.H{"token": token})
}
