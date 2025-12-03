package db

import (
	"fmt"
	"go-sosmed-app/internal/models"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

	fmt.Println("Connected to database succsessfully")

	db.AutoMigrate(&models.User{}, &models.Post{}, &models.Comment{}, &models.Like{}, &models.PostMedia{})

	DB = db
}
