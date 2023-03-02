package services

import (
	"go-echo/helpers"
	"go-echo/locales"
	"go-echo/models"
	"go-echo/queries"
	"go-echo/structs"

	"gorm.io/gorm"
)

type AuthService struct {
	DB *gorm.DB
}

func (as *AuthService) Show(id uint) (models.UserRoleWithPermission, int, error) {
	var user models.UserRoleWithPermission
	query := as.DB.Scopes(queries.RoleUserWithPermissionPreload()).First(&user, id)

	statusCode, err := helpers.ErrorDatabaseDynamic(query.Error, helpers.DatabaseDynamicMessage{
		NotFound: locales.LocalesGet("user.validation.notFound"),
	})

	return user, statusCode, err
}

func (as *AuthService) Store(payload structs.UserStoreForm) (*models.UserResponse, error) {
	hashedPassword, err := helpers.PasswordHash(payload.Password)
	if err != nil {
		return nil, err
	}

	user := models.UserResponse{
		Name:     payload.Name,
		Email:    payload.Email,
		Password: hashedPassword,
	}
	if err := as.DB.Create(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
