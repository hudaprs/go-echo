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

	query := us.DB.Scopes(helpers.Paginate(users, &pagination, us.DB)).Scopes(queries.RoleUserPreload()).Find(&users)
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

	// Assign payload.ID to user.id
	payload.ID = &user.ID

	// Check if theres any roles from payload
	// If no, let var roles blank
	if len(payload.Roles) > 0 {
		if err := us.DB.Where(payload.Roles).Find(&roles); err.Error != nil {
			return user, err.Error
		}
	}

	// Force to clear previous roles
	// Ignore when creating user
	if query := us.DB.Where(models.RoleUser{UserID: *payload.ID}).Delete(&models.RoleUser{}); query.Error != nil {
		return user, query.Error
	}

	// Force to create new roles after create / update
	var assignedUserRoles []models.RoleUserResponse
	if len(roles) > 0 {
		for _, role := range roles {
			newUserRole := &models.RoleUserResponse{
				RoleID: role.ID,
				UserID: *payload.ID,
				Name:   role.Name,
			}

			if query := us.DB.Create(newUserRole); query.Error != nil {
				return user, query.Error
			}

			// Map new data id to be role id
			newUserRole.ID = role.ID

			assignedUserRoles = append(assignedUserRoles, *newUserRole)
		}
	}

	return models.UserWithRoleResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Roles:     assignedUserRoles,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (us *UserService) Show(payload structs.UserAttrsFind) (models.UserWithRoleResponse, int, error) {
	var user models.UserWithRoleResponse

	query := us.DB.Scopes(queries.RoleUserPreload()).Where("id = ?", payload.ID).Or("name = ?", payload.Name).Or("email = ?", payload.Email).First(&user)
	queryStatusCode := helpers.ValidateNotFoundData(query.Error)

	return user, queryStatusCode, query.Error
}

func (us *UserService) Delete(payload structs.UserAttrsFind) (models.UserWithRoleResponse, int, error) {
	var user models.UserWithRoleResponse

	query := us.DB.Scopes(queries.RoleUserPreload()).Where("id = ?", payload.ID).Or("name = ?", payload.Name).Or("email = ?", payload.Email).First(&user).Delete(user)
	queryStatusCode := helpers.ValidateNotFoundData(query.Error)

	return user, queryStatusCode, query.Error
}

func (us *UserService) CheckEmail(email string) (bool, int, models.UserResponse, error) {
	var user models.UserResponse
	query := us.DB.Where("email = ?", email).First(&user)

	findUserStatusCode := helpers.ValidateNotFoundData(query.Error)

	return findUserStatusCode == 200, findUserStatusCode, user, query.Error
}
