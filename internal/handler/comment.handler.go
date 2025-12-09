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

func validateComment(comment string) (bool, string, string) {
	trimedComment := strings.TrimSpace(comment)

	if trimedComment == "" {
		return false, "", "Comment cannot be empty"
	}

	if len(trimedComment) > 500 {
		return false, "", "comment_text exceed 500 character after space trim"
	}

	return true, trimedComment, ""
}

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

	isValid, commentText, message := validateComment(req.CommentText)

	if !isValid {
		c.JSON(400, gin.H{
			"success": false,
			"message": message,
		})
		return
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

func UpdateComment(c *gin.Context) {
	userId := c.GetUint("userId")

	commentIdString := c.Param("id")

	commentId, err := strconv.ParseUint(commentIdString, 10, 64)

	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Invalid id format",
		})
		return
	}

	var comment models.Comment

	if err := db.DB.First(&comment, commentId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(404, gin.H{
				"success": false,
				"message": "Comment not found",
			})
			return
		}

		c.JSON(500, gin.H{
			"success": false,
			"message": "Failed to fetch comment: " + err.Error(),
		})
		return
	}

	if comment.UserID != userId {
		c.JSON(403, gin.H{
			"success": false,
			"message": "User not allowed",
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

	isValid, commentText, message := validateComment(req.CommentText)

	if !isValid {
		c.JSON(400, gin.H{
			"success": false,
			"message": message,
		})
		return
	}

	if err := db.DB.Model(&comment).Where("id = ?", commentId).Update("comment_text", commentText).Error; err != nil {
		c.JSON(500, gin.H{
			"success": false,
			"message": "Failed to update comment: " + err.Error(),
		})
	}

	c.JSON(200, gin.H{
		"success": true,
		"data":    comment,
	})
}

func DeleteComment(c *gin.Context) {
	userId := c.GetUint("userId")

	commentIdString := c.Param("id")

	commentId, err := strconv.ParseUint(commentIdString, 10, 64)

	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Invalid comment id format",
		})
		return
	}

	var comment models.Comment

	if err := db.DB.First(&comment, commentId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(404, gin.H{
				"success": false,
				"message": "Comment not found",
			})
			return
		}

		c.JSON(500, gin.H{
			"success": false,
			"message": "Failed to check if comment exists: " + err.Error(),
		})
		return
	}

	if comment.UserID != userId {
		c.JSON(403, gin.H{
			"success": false,
			"message": "User not allowed",
		})
		return
	}

	if err := db.DB.Delete(&comment).Error; err != nil {
		c.JSON(500, gin.H{
			"success": false,
			"message": "Failed to delete comment: " + err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"data":    comment,
	})
}
