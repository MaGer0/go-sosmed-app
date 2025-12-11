package handler

import (
	"errors"
	"fmt"
	"go-sosmed-app/internal/db"
	"go-sosmed-app/internal/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CountPostLikes(c *gin.Context) {
	stringPostId := c.Param("postId")

	postId, err := strconv.ParseUint(stringPostId, 10, 64)

	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Invalid PostID format",
		})
		return
	}

	var postExists int64
	if err := db.DB.Model(&models.Post{}).Where("id = ?", postId).Count(&postExists).Error; err != nil {
		c.JSON(500, gin.H{
			"success": false,
			"message": fmt.Sprintf("Failed to check if Post with ID %s exists", stringPostId),
		})
		return
	}

	if postExists == 0 {
		c.JSON(404, gin.H{
			"success": false,
			"message": fmt.Sprintf("Post with ID %s not exists", stringPostId),
		})
		return
	}

	var likesCount int64
	if err := db.DB.Model(&models.Like{}).Where("post_id = ?", postId).Count(&likesCount).Error; err != nil {
		c.JSON(500, gin.H{
			"success": false,
			"message": fmt.Sprintf("Failed to fetch Like with PostID %s", stringPostId),
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"data": gin.H{
			"likes_count": likesCount,
		},
	})
}

func LikePost(c *gin.Context) {
	userId := c.GetUint("userid")
	stringPostId := c.Param("postId")

	postId, err := strconv.ParseUint(stringPostId, 10, 64)

	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"data":    "Invalid ID format",
		})
		return
	}

	var post models.Post
	if err := db.DB.Select("id").First(&post, postId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(404, gin.H{
				"success": false,
				"data":    fmt.Sprintf("Post with ID %s not found", stringPostId),
			})
			return
		}

		c.JSON(500, gin.H{
			"success": false,
			"message": "Failed to check if Post exists",
		})
		return
	}

	var like models.Like
	err = db.DB.Where("user_id = ? AND post_id = ?", userId, post.ID).First(&like).Error

	if err == nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Like is already exists",
		})
		return
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(500, gin.H{
			"success": false,
			"message": fmt.Sprintf("Failed to check if Like exsits: %s", err.Error()),
		})
		return
	}

	like = models.Like{
		UserID: userId,
		PostID: post.ID,
	}

	if err := db.DB.Create(&like).Error; err != nil {
		c.JSON(500, gin.H{
			"success": false,
			"message": "Failed to create Like",
		})
		return
	}

	c.JSON(201, gin.H{
		"success": true,
		"data":    like,
	})
}
