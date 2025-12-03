package models

import "gorm.io/gorm"

type PostMedia struct {
	gorm.Model
	PostID     uint   `json:"post_id"`
	ImageURL   string `json:"image_url"`
	OrderIndex int    `json:"order_index"`
}