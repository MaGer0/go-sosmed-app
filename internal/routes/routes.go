package routes

import (
	"go-sosmed-app/internal/handler"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")

	// Auth
	api.POST("/register", handler.Register)
	api.POST("/login", handler.Login)

	// Post
	api.GET("/posts", handler.AuthMiddleware, handler.GetPosts)
	api.POST("/posts", handler.AuthMiddleware, handler.CreatePost)
	api.PATCH("/posts/:id", handler.AuthMiddleware, handler.UpdatePostCaption)
	api.DELETE("/posts/:id", handler.AuthMiddleware, handler.DeletePost)

	//Comment
	api.POST("/comments/:postId", handler.AuthMiddleware, handler.AddComment)
	api.PATCH("/comments/:id", handler.AuthMiddleware, handler.UpdateComment)
	api.DELETE("/comments/:id", handler.AuthMiddleware, handler.DeleteComment)

	// Post Media
	api.GET("/post-media/:postId", handler.AuthMiddleware, handler.GetPostMedia)
	api.PATCH("/post-media/:postId", handler.AuthMiddleware, handler.UpdatePostMedia)
	api.DELETE("/post-media/:postId", handler.AuthMiddleware, handler.DeletePostMedia)
}
