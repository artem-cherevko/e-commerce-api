package database

import (
	"e-commerce-api/internal/modules/models"
	"log"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(&models.User{}, &models.UserSessions{}, &models.Product{})
	if err != nil {
		log.Fatalf("error while trying to autoMigrate: %s", err.Error())
	}
}
