package models

import (
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	UserID  uint   `json:"user_id"`
	Likes   []Like `json:"likes"`
	Caption string `json:"caption"`
}
