package migrations

import (
	"fmt"
	"go-echo/models"

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
		&models.RoleUser{},
		&models.RolePermission{},
	)
	if err != nil {
		panic("Migration: something went wrong when start to migrate " + err.Error())
	}

	fmt.Println("Migration: successfully migrate models to database")
}
