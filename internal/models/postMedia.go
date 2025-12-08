package models

import "time"

type PostMedia struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	PostID     uint      `json:"post_id"`
	ImageURL   string    `json:"image_url"`
	OrderIndex int       `json:"order_index"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
