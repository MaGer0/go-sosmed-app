package models

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	PostID      uint   `json:"post_id"`
	UserID      uint   `json:"user_id"`
	CommentText string `json:"comment_text"`
}
