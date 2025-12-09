package handler

import (
	"errors"
	"fmt"
	"go-sosmed-app/internal/db"
	"go-sosmed-app/internal/models"
	"go-sosmed-app/internal/utils"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Perbaikin, ini mah itungannya belum mandiri (masih dipanggil di function lain, belum benar-benar handler sendiri)
func UploadMedia(c *gin.Context, postId uint, tx *gorm.DB) ([]models.PostMedia, error) {
	form, err := c.MultipartForm()

	if err != nil {
		return nil, fmt.Errorf("failed to parse multipart form: %w", err)
	}

	files := form.File["files"]

	if len(files) == 0 {
		return nil, nil
	}

	var postMedia []models.PostMedia

	for _, file := range files {
		contentType := file.Header.Get("Content-Type")

		if !strings.HasPrefix(contentType, "image/") && !strings.HasPrefix(contentType, "video/") {
			return nil, fmt.Errorf("invalid file type: %s", file.Filename)
		}

		uniqueName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), file.Filename)
		relativePath := "uploads/" + uniqueName
		fullPath := "./storage/uploads/" + uniqueName

		if err := c.SaveUploadedFile(file, fullPath); err != nil {
			return nil, fmt.Errorf("failed to save file: %w", err)
		}

		postMedia = append(postMedia, models.PostMedia{
			PostID:   postId,
			ImageURL: relativePath,
		})
	}

	if len(postMedia) > 0 {
		if err := tx.Create(&postMedia).Error; err != nil {
			return nil, fmt.Errorf("failed to create PostMedia: %w", err)
		}
	}

	return postMedia, nil
}

func UpdatePostMedia(c *gin.Context) {
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

	if err := db.DB.Select("id", "user_id").Where("id = ?", postId).First(&post).Error; err != nil {
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

	if post.UserID != userId {
		c.JSON(403, gin.H{
			"success": false,
			"message": "User not allowed",
		})
		return
	}

	var req struct {
		MediaIdsToDelete []uint `form:"media_ids_to_delete"`
	}

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": utils.ParseErrorMessage(err),
		})
		return
	}

	form, err := c.MultipartForm()

	if err != nil {
		c.JSON(500, gin.H{
			"success": false,
			"message": "Failed to parse multipart form: " + err.Error(),
		})
		return
	}

	files := form.File["files"]

	var postMedia []models.PostMedia

	for _, file := range files {
		contentType := file.Header.Get("Content-Type")

		if !strings.HasPrefix(contentType, "image/") && !strings.HasPrefix(contentType, "video/") {
			c.JSON(400, gin.H{
				"success": false,
				"message": "Invalid file type: " + file.Filename,
			})
			return
		}

		uniqueName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), file.Filename)
		relativePath := "uploads/" + uniqueName
		fullPath := "./storage/uploads/" + uniqueName

		if err := c.SaveUploadedFile(file, fullPath); err != nil {
			c.JSON(500, gin.H{
				"success": false,
				"message": "Failed to save file: " + err.Error(),
			})
			return
		}

		postMedia = append(postMedia, models.PostMedia{
			PostID:   uint(postId),
			ImageURL: relativePath,
		})
	}

	if err := db.DB.Create(&postMedia).Error; err != nil {
		c.JSON(500, gin.H{
			"success": false,
			"message": "Failed to save post media: " + err.Error(),
		})
		return
	}

	var mediaToDelete []models.PostMedia
	tes := 0
	if len(req.MediaIdsToDelete) > 0 {
		tes = len(req.MediaIdsToDelete)

		if err := db.DB.Where("id IN ?", req.MediaIdsToDelete).Where("post_id = ?", postId).Find(&mediaToDelete).Error; err != nil {
			c.JSON(500, gin.H{
				"success": false,
				"message": "Failed to fetch post media: " + err.Error(),
			})
			return
		}

		for _, mediumToDelete := range mediaToDelete {
			os.Remove(filepath.Join("storage", mediumToDelete.ImageURL))
		}

		if err := db.DB.Delete(&mediaToDelete).Error; err != nil {
			c.JSON(500, gin.H{
				"success": false,
				"message": "Failed to delete post media: " + err.Error(),
			})
			return
		}
	}

	c.JSON(200, gin.H{
		"success": true,
		"tes":     tes,
	})
}

// func UpdatePostMedia(c *gin.Context) {
// 	postIdString := c.Param("postId")

// 	postId, err := strconv.ParseUint(postIdString)

// 	if err != nil {

// 	}
// }
