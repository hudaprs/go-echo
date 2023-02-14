package services

import (
	"echo-rest/helpers"
	"echo-rest/models"
	"echo-rest/queries"
	"echo-rest/structs"

	"gorm.io/gorm"
)

type AuthService struct {
	DB *gorm.DB
}

func (as *AuthService) Show(id uint) (models.UserRoleWithPermission, int, error) {
	var user models.UserRoleWithPermission
	query := as.DB.Scopes(queries.RoleUserWithPermissionPreload()).First(&user, id)

	statusCode, err := helpers.ErrorDatabaseNotFound(query.Error)

	return user, statusCode, err
}

func (as *AuthService) Store(payload structs.UserStoreForm) (models.UserResponse, error) {
	hashedPassword, err := helpers.PasswordHash(payload.Password)

	if err != nil {
		panic("User Store: failed when start to hash password")
	}

	user := models.UserResponse{
		Name:     payload.Name,
		Email:    payload.Email,
		Password: hashedPassword,
	}

	query := as.DB.Create(&user)

	return user, query.Error
}
