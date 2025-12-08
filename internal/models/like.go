package models

import "time"

type Like struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"uniqueIndex:idx_user_post_uniq_together"`
	PostID    uint      `json:"post_id" gorm:"uniqueIndex:idx_user_post_uniq_together"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
