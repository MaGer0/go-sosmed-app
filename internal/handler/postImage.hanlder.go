package handler

import (
	"fmt"
	"go-sosmed-app/internal/models"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UploadImage(c *gin.Context, postId uint, tx *gorm.DB) ([]models.PostMedia, error) {
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
