package migrations

import (
	"echo-rest/database"
	"echo-rest/models"
	"fmt"
	"os"

	"gorm.io/gorm"
)

func InitMigration(db *gorm.DB) {
	err := database.DatabaseConnection().AutoMigrate(&models.Todo{}, &models.User{}, &models.RefreshToken{})

	if err != nil {
		fmt.Println("Migration: something went wrong when start to migrate", err)
		os.Exit(0)
	}

	fmt.Println("Migration: successfully migrate models to database")
}
