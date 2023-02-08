package services

import (
	"echo-rest/helpers"
	"echo-rest/models"
	"echo-rest/queries"
	"echo-rest/structs"

	"gorm.io/gorm"
)

type UserService struct {
	DB *gorm.DB
}

func (us *UserService) Index(pagination helpers.Pagination) (*helpers.Pagination, error) {
	var users []models.UserWithRoleResponse

	query := us.DB.Scopes(helpers.Paginate(users, &pagination, us.DB)).Preload("Roles", queries.RolesMap).Find(&users)
	pagination.Rows = users

	return &pagination, query.Error
}

func (us *UserService) Store(payload structs.UserCreateForm) (models.UserWithRoleResponse, error) {
	var user models.UserWithRoleResponse
	var roles []models.RoleResponse

	// Make default has password for the first time creating an user
	hashedPassword, err := helpers.PasswordHash("password")
	if err != nil {
		return models.UserWithRoleResponse{}, err
	}

	// Create new user
	user.Name = payload.Name
	user.Email = payload.Email
	user.Password = hashedPassword
	if err := us.DB.Save(&user); err.Error != nil {
		return user, err.Error
	}

	// Load roles
	if err := us.DB.Find(&roles, payload.Roles); err.Error != nil {
		return user, err.Error
	}

	// Assign role to newest user
	if err := us.DB.Model(&user).Association("Roles").Append(roles); err != nil {
		return user, err
	}

	return user, nil
}

func (us *UserService) Show(payload structs.UserAttrsFind) (models.UserWithRoleResponse, int, error) {
	var user models.UserWithRoleResponse

	query := us.DB.Preload("Roles").Where("id = ?", payload.ID).Or("name = ?", payload.Name).Or("email = ?", payload.Email).First(&user)
	queryStatusCode := helpers.ValidateNotFoundData(query.Error)

	return user, queryStatusCode, query.Error
}

func (us *UserService) Update() {
	//
}

func (us *UserService) Delete(payload structs.UserAttrsFind) (models.UserWithRoleResponse, int, error) {
	var user models.UserWithRoleResponse

	query := us.DB.Preload("Roles").Where("id = ?", payload.ID).Or("name = ?", payload.Name).Or("email = ?", payload.Email).First(&user).Delete(user)
	queryStatusCode := helpers.ValidateNotFoundData(query.Error)

	return user, queryStatusCode, query.Error
}

func (us *UserService) CheckEmail(email string) (bool, int, models.UserResponse, error) {
	var user models.UserResponse
	query := us.DB.Where("email = ?", email).First(&user)

	findUserStatusCode := helpers.ValidateNotFoundData(query.Error)

	return findUserStatusCode == 200, findUserStatusCode, user, query.Error
}
