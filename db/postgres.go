package db

import (
	"log"

	"gorm.io/gorm"
	"gorm.io/driver/postgres"

	"payment-service/internal/model"
)

func Connect(strConn string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(strConn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err = db.Exec("SELECT 1").Error; err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	if err = db.AutoMigrate(&model.Transfer{}, &model.Account{}); err != nil {
		log.Fatal("Failed to auto-migrate database:", err)
	}

	return db
}