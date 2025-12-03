package controllers

import (
	"go-sosmed-app/internal/db"
	"go-sosmed-app/internal/models"

	"github.com/gin-gonic/gin"
)

func GetPosts(c *gin.Context) {
	userId := c.GetUint("userId")
	var posts []models.Post

	if err := db.DB.Where("user_id = ?", userId).Find(&posts).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, posts)
}

func CreatePost(c *gin.Context) {
	userId := c.GetUint("userId")
	var req struct {
		Caption string `form:"caption" binding:"required"`
	}

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	post := models.Post{
		UserID:  userId,
		Caption: req.Caption,
	}

	if err := db.DB.Create(&post).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	form, _ := c.MultipartForm()
	files := form.File["files"]

	var postMedia []models.PostMedia

	for _, file := range files {
		dst := "./storage/uploads/" + file.Filename
		c.SaveUploadedFile(file, dst)

		postMedia = append(postMedia, models.PostMedia{
			PostID:   post.ID,
			ImageURL: dst,
		})
	}

	if len(postMedia) > 0 {
		if err := db.DB.Create(&postMedia).Error; err != nil {
			c.JSON(500, err.Error())
			return
		}
	}

	c.JSON(201, gin.H{
		"message": "post created successfully",
		"post":    post,
		"media":   postMedia,
	})

}
