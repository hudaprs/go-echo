package services

import (
	"echo-rest/helpers"
	"echo-rest/models"
	"echo-rest/structs"

	"gorm.io/gorm"
)

type UserService struct {
	DB *gorm.DB
}

func (us *UserService) Index(pagination helpers.Pagination) (*helpers.Pagination, error) {
	var users []models.UserWithRoleResponse

	query := us.DB.Scopes(helpers.Paginate(users, &pagination, us.DB)).Preload("Roles").Find(&users)
	pagination.Rows = users

	return &pagination, query.Error
}

func (us *UserService) StoreOrUpdate(payload structs.UserCreateEditForm) (models.UserWithRoleResponse, error) {
	var user models.UserWithRoleResponse
	var roles []models.RoleResponse

	// Make default has password for the first time creating an user
	hashedPassword, err := helpers.PasswordHash("password")
	if err != nil {
		return models.UserWithRoleResponse{}, err
	}

	// Assign value
	if payload.ID != nil {
		query := us.DB.Find(&user, payload.ID)
		if query.Error != nil {
			return user, query.Error
		}
	}
	user.Name = payload.Name
	user.Email = payload.Email
	if payload.ID == nil {
		user.Password = hashedPassword
	}

	// Create / Update new user
	if err := us.DB.Save(&user); err.Error != nil {
		return user, err.Error
	}

	// Check if theres any roles from payload
	// If no, let var roles blank
	if len(payload.Roles) > 0 {
		if err := us.DB.Where(payload.Roles).Find(&roles); err.Error != nil {
			return user, err.Error
		}
	}

	// Force to clear roles associated to users
	if err := us.DB.Model(&user).Association("Roles").Clear(); err != nil {
		return user, err
	}

	// Force to create new roles after create / update
	if len(roles) > 0 {
		if err := us.DB.Model(&user).Association("Roles").Append(roles); err != nil {
			return user, err
		}
	}

	return user, nil
}

func (us *UserService) Show(payload structs.UserAttrsFind) (models.UserWithRoleResponse, int, error) {
	var user models.UserWithRoleResponse

	query := us.DB.Preload("Roles").Where("id = ?", payload.ID).Or("name = ?", payload.Name).Or("email = ?", payload.Email).First(&user)
	queryStatusCode := helpers.ValidateNotFoundData(query.Error)

	return user, queryStatusCode, query.Error
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
