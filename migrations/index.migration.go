package migrations

import (
	"echo-rest/models"
	"fmt"

	"gorm.io/gorm"
)

func InitMigration(db *gorm.DB) {
	// Migration List
	err := db.AutoMigrate(
		&models.Todo{},
		&models.User{},
		&models.RefreshToken{},
		&models.Role{},
		&models.Permission{},
	)
	if err != nil {
		panic("Migration: something went wrong when start to migrate" + err.Error())
	}

	fmt.Println("Migration: successfully migrate models to database")
}
