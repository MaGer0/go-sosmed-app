package handler

import (
	"errors"
	"go-sosmed-app/internal/db"
	"go-sosmed-app/internal/models"
	"go-sosmed-app/internal/utils"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AddComment(c *gin.Context) {
	userId := c.GetUint("userId")
	postIdString := c.Param("postId")

	postId, err := strconv.ParseUint(postIdString, 10, 64)

	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Invalid post id format",
		})
		return
	}

	var post models.Post

	if err := db.DB.Select("id").First(&post, postId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(404, gin.H{
				"success": false,
				"message": "Post not found",
			})
			return
		}

		c.JSON(500, gin.H{
			"success": false,
			"message": "Failed to check if post exists: " + err.Error(),
		})
		return
	}

	var req struct {
		CommentText string `json:"comment_text" binding:"required,max=500"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": utils.ParseErrorMessage(err),
		})
		return
	}

	commentText := strings.TrimSpace(req.CommentText)

	if commentText == "" {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Comment cannot be empty",
		})
		return
	}

	if len(commentText) > 500 {
		c.JSON(400, gin.H{
			"success": false,
			"message": "comment_text exceed 500 character after space trim",
		})
	}

	comment := models.Comment{
		UserID:      userId,
		PostID:      post.ID,
		CommentText: commentText,
	}

	if err := db.DB.Create(&comment).Error; err != nil {
		c.JSON(500, gin.H{
			"success": false,
			"message": "Failed to add comment: " + err.Error(),
		})
		return
	}

	c.JSON(201, gin.H{
		"success": true,
		"data":    comment,
	})
}
