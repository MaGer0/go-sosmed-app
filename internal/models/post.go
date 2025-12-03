package models

import (
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	UserID  uint   `json:"user_id"`
	Caption string `json:"caption"`
	Likes   []Like `json:"likes"`
}
