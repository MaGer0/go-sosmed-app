package handler

import (
	"errors"
	"go-sosmed-app/internal/db"
	"go-sosmed-app/internal/models"
	"go-sosmed-app/internal/utils"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func validateCaption(caption string) (string, bool, string) {
	trimedCaption := strings.TrimSpace(caption)

	if trimedCaption == "" {
		return trimedCaption, false, "Field caption cannot be empty"
	}

	if len(trimedCaption) > 500 {
		return trimedCaption, false, "Field caption cannot exceed 500 character"
	}

	return trimedCaption, true, ""
}

func GetPosts(c *gin.Context) {
	userId := c.GetUint("userId")
	var posts []models.Post

	if err := db.DB.Where("user_id = ?", userId).Find(&posts).Error; err != nil {
		c.JSON(500, gin.H{
			"success": false,
			"message": "Failed to fetch posts: " + err.Error(),
		})
		return
	}

	c.JSON(200, posts)
}

func CreatePost(c *gin.Context) {
	userId := c.GetUint("userId")
	var req struct {
		Caption string `form:"caption" binding:"required,max=500"`
	}

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": utils.ParseErrorMessage(err),
		})
		return
	}

	trimedCaption, isValid, errCaptionValidate := validateCaption(req.Caption)

	if !isValid {
		c.JSON(400, gin.H{
			"success": false,
			"message": errCaptionValidate,
		})
	}

	post := models.Post{
		UserID:  userId,
		Caption: trimedCaption,
	}

	tx := db.DB.Begin()

	if err := tx.Create(&post).Error; err != nil {
		tx.Rollback()
		c.JSON(500, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	postMedia, err := UploadMedia(c, post.ID, tx)

	if err != nil {
		tx.Rollback()
		c.JSON(400, gin.H{
			"success": false,
			"message": "Failed to upload media: " + err.Error(),
		})
		return
	}

	tx.Commit()

	c.JSON(201, gin.H{
		"success": true,
		"message": "post created successfully",
		"data": gin.H{
			"post":       post,
			"post_media": postMedia,
		},
	})

}

func UpdatePostCaption(c *gin.Context) {
	userId := c.GetUint("userId")
	idString := c.Param("id")

	id, err := strconv.ParseUint(idString, 10, 64)

	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Invalid id format",
		})
		return
	}

	var post models.Post

	if err := db.DB.First(&post, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(404, gin.H{
				"success": false,
				"message": "Post not found",
			})
			return
		}

		c.JSON(500, gin.H{
			"success": false,
			"message": "Failed to fetch post: " + err.Error()})
		return
	}

	if userId != post.UserID {
		c.JSON(403, gin.H{
			"success": false,
			"message": "User not allowed",
		})
		return
	}

	var req struct {
		Caption string `json:"caption" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": utils.ParseErrorMessage(err),
		})
		return
	}

	if err := db.DB.Model(&post).Update("caption", req.Caption).Error; err != nil {
		c.JSON(500, gin.H{
			"success": false,
			"message": "Failed to update caption: " + err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"data":    post,
	})
}

func DeletePost(c *gin.Context) {
	userId := c.GetUint("userId")
	idString := c.Param("id")

	id, err := strconv.ParseUint(idString, 10, 64)

	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Invalid id format",
		})
		return
	}

	var post models.Post

	if err := db.DB.First(&post, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(404, gin.H{
				"success": false,
				"message": "Post not found",
			})
			return
		}

		c.JSON(500, gin.H{
			"success": false,
			"message": "failed to fetch post: " + err.Error(),
		})
	}

	if post.UserID != userId {
		c.JSON(403, gin.H{
			"success": false,
			"message": "User not allowed",
		})
		return
	}

	// Delete postMedia
	var postMedia []models.PostMedia

	if err := db.DB.Where("post_id = ?", post.ID).Find(&postMedia).Error; err != nil {
		c.JSON(500, gin.H{
			"success": false,
			"message": "Failed to fetch post media: " + err.Error(),
		})
		return
	}

	for _, postMedium := range postMedia {
		if err := os.Remove(filepath.Join("storage", postMedium.ImageURL)); err != nil {
			c.JSON(500, gin.H{
				"success": false,
				"message": "Failed to delete post media: " + err.Error(),
			})
			return
		}
	}

	tx := db.DB.Begin()

	if err := tx.Where("post_id = ?", post.ID).Delete(&models.PostMedia{}).Error; err != nil {
		tx.Rollback()
		c.JSON(500, gin.H{
			"success": false,
			"message": "Failed to delete post media: " + err.Error(),
		})
		return
	}

	if err := tx.Delete(&post).Error; err != nil {
		tx.Rollback()
		c.JSON(500, gin.H{
			"success": false,
			"message": "Failed to delete post: " + err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"data":    post,
	})
}
