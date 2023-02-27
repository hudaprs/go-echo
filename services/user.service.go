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

func (us *UserService) createRoleUser(roleIds []uint, userId uint) ([]models.RoleUserResponse, error) {
	roleUserResponse := []models.RoleUserResponse{}

	// Check if user didn't include payload
	if len(roleIds) == 0 {
		if query := us.DB.Where("user_id = ?", userId).Delete(&models.RoleUserResponse{}); query.Error != nil {
			return roleUserResponse, query.Error
		}
	}

	// Check if user already have role before
	// If user didn't sent the previous role, remove, but assign that new one
	if query := us.DB.Where("user_id = ?", userId).Where("role_id NOT IN ?", roleIds).Delete(&models.RoleUserResponse{}); query.Error != nil {
		return roleUserResponse, query.Error
	}

	// Check if user have existing role
	var existedUserRoles []models.RoleUser
	if query := us.DB.Where(models.RoleUser{UserID: userId}).Where("role_id IN ?", roleIds).Find(&existedUserRoles); query.Error != nil {
		return roleUserResponse, query.Error
	}

	// Find unique role that never be assigned before to user
	var assignedRoleIds []uint
	for _, roleId := range roleIds {
		// Make identifier to skip if not exists
		skip := false
		for _, existedUserRole := range existedUserRoles {
			// If data found
			// Just don't do anything
			if roleId == existedUserRole.RoleID {
				// If data match, force to true, don't do anything
				skip = true
				break
			}
		}

		// If role not found, just make a new one
		if !skip {
			assignedRoleIds = append(assignedRoleIds, roleId)
		}
	}

	// Check if theres any unique role that user didn't assign before
	// Create that unique roles and assign to specific users
	if len(assignedRoleIds) > 0 {
		for _, roleId := range assignedRoleIds {
			// Check if role exists
			if query := us.DB.First(&models.Role{}, roleId); query.Error != nil {
				return roleUserResponse, query.Error
			}

			newUserRole := &models.RoleUserResponse{
				RoleID: roleId,
				UserID: userId,
			}

			if query := us.DB.Create(newUserRole); query.Error != nil {
				return roleUserResponse, query.Error
			}
		}
	}

	// Look up new assigned role through database
	if query := us.DB.Select("role_users.*, roles.name, roles.id").Where("user_id = ?", userId).Joins("left join roles ON roles.id = role_users.role_id").Order("created_at desc").Find(&roleUserResponse); query.Error != nil {
		return roleUserResponse, query.Error
	}

	return roleUserResponse, nil
}

func (us *UserService) Index(pagination helpers.Pagination) (*helpers.Pagination, error) {
	var users []models.UserResponse

	query := us.DB.Scopes(helpers.Paginate(users, &pagination, us.DB)).Find(&users)
	pagination.Rows = users

	return &pagination, query.Error
}

func (us *UserService) StoreOrUpdate(payload structs.UserCreateEditForm, isCreate bool) (models.UserWithRoleResponse, error) {
	var user models.UserWithRoleResponse
	var assignedUserRoles []models.RoleUserResponse

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

	// Create or Update user role
	userRoles, err := us.createRoleUser(payload.Roles, *payload.ID)
	if err != nil {
		return user, err
	}
	assignedUserRoles = userRoles

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
