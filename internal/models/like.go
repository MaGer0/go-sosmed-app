package models

import "gorm.io/gorm"

type Like struct {
	gorm.Model
	UserID uint `json:"user_id" gorm:"uniqueIndex:idx_user_post_uniq_together"`
	PostID uint `json:"post_id" gorm:"uniqueIndex:idx_user_post_uniq_together"`
}
