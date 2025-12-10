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
	}

	for _, file := range files {
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

	if len(req.MediaIdsToDelete) > 0 {
		var mediaToDelete []models.PostMedia
		if err := db.DB.Where("id IN ?", req.MediaIdsToDelete).Where("post_id = ?", postId).Find(&mediaToDelete).Error; err != nil {
			c.JSON(500, gin.H{
				"success": false,
				"message": "Failed to fetch post media: " + err.Error(),
			})
			return
		}

		if len(mediaToDelete) == 0 {
			c.JSON(404, gin.H{
				"success": false,
				"message": "Post media cannot be found",
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
				"data":    mediaToDelete,
			})
			return
		}
	}

	if err := db.DB.Where("post_id = ?", postId).Find(&postMedia).Error; err != nil {
		c.JSON(500, gin.H{
			"success": false,
			"message": "Failed to fetch post media: " + err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"data":    postMedia,
	})
}

func DeletePostMedia(c *gin.Context) {
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
	if err := db.DB.Select("user_id").First(&post, postId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(404, gin.H{
				"success": false,
				"message": "Post not found",
			})
			return
		}

		c.JSON(500, gin.H{
			"success": false,
			"message": "Failed to fetch post: " + err.Error(),
		})
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
		PostMediaIdsToDelete []uint `json:"post_media_ids_to_delete" bind:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": utils.ParseErrorMessage(err),
		})
		return
	}

	if len(req.PostMediaIdsToDelete) == 0 {
		c.JSON(400, gin.H{
			"success": false,
			"message": "post_media_ids_to_delete is empty",
		})
		return
	}

	var PostMediaIdsToDelete []models.PostMedia
	if err := db.DB.Select("id").Where("id IN ?", req.PostMediaIdsToDelete).Find(&PostMediaIdsToDelete).Error; err != nil {
		c.JSON(500, gin.H{
			"success": false,
			"message": "Failed to check if ID exists: " + err.Error(),
		})
		return
	}

	var stringIdsSlice []string

	for _, postMediumId := range req.PostMediaIdsToDelete {
		stringIdsSlice = append(stringIdsSlice, fmt.Sprintf("%d", postMediumId))
	}

	stringIds := strings.Join(stringIdsSlice, ", ")

	if len(PostMediaIdsToDelete) == 0 {
		c.JSON(404, gin.H{
			"success": false,
			"message": "Post media with ID " + stringIds + " not found",
		})
		return
	}

	if err := db.DB.Where("id IN ?", req.PostMediaIdsToDelete).Where("post_id = ?", postId).Delete(&models.PostMedia{}).Error; err != nil {
		c.JSON(500, gin.H{
			"success": false,
			"message": "Failed to delete post media: " + err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"message": "Post media with ID " + stringIds + " deleted successfully",
	})

}
