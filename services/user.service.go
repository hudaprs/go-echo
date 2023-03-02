package services

import (
	"go-echo/database"
	"go-echo/helpers"
	"go-echo/models"
	"go-echo/queries"
	"go-echo/structs"

	"gorm.io/gorm"
)

type UserService struct {
	DB *gorm.DB
}

func syncRoleUser(tx *gorm.DB, roleIds []uint, userId uint) ([]models.RoleUserResponse, error) {
	roleUserResponse := []models.RoleUserResponse{}

	// Check if user didn't include payload
	if len(roleIds) == 0 {
		if err := tx.Where("user_id = ?", userId).Delete(&models.RoleUserResponse{}).Error; err != nil {
			return roleUserResponse, err
		}
	}

	// Check if user already have role before
	// If user didn't sent the previous role, remove, but assign that new one
	if query := tx.Where("user_id = ?", userId).Where("role_id NOT IN ?", roleIds).Delete(&models.RoleUserResponse{}); query.Error != nil {
		return roleUserResponse, query.Error
	}

	// Check if user have existing role
	var existedUserRoles []models.RoleUser
	if query := tx.Where(models.RoleUser{UserID: userId}).Where("role_id IN ?", roleIds).Find(&existedUserRoles); query.Error != nil {
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
			if query := tx.First(&models.Role{}, roleId); query.Error != nil {
				return roleUserResponse, query.Error
			}

			newUserRole := &models.RoleUserResponse{
				RoleID: roleId,
				UserID: userId,
			}

			if query := tx.Create(newUserRole); query.Error != nil {
				return roleUserResponse, query.Error
			}
		}
	}

	// Look up new assigned role through database
	if query := tx.Select("role_users.*, roles.name, roles.id").Where("user_id = ?", userId).Joins("left join roles ON roles.id = role_users.role_id").Order("created_at desc").Find(&roleUserResponse); query.Error != nil {
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

func (us *UserService) StoreOrUpdate(payload structs.UserCreateEditForm) (*models.UserWithRoleResponse, error) {
	user := &models.UserWithRoleResponse{}
	tx, err := database.BeginTransaction(us.DB)
	if err != nil {
		return nil, err
	}

	// Make default has password for the first time creating an user
	hashedPassword, err := helpers.PasswordHash("password")
	if err != nil {
		return nil, err
	}

	// Assign value
	if payload.ID != nil {
		query := tx.Find(&user, payload.ID)
		if query.Error != nil {
			return nil, query.Error
		}
	}
	user.Name = payload.Name
	user.Email = payload.Email
	if payload.ID == nil {
		user.Password = hashedPassword
	}

	// Create / Update new user
	if err := tx.Save(&user); err.Error != nil {
		return nil, err.Error
	}

	// Assign payload.ID to user.id
	payload.ID = &user.ID

	// Create or Update user role
	userRoles, err := syncRoleUser(tx, payload.Roles, *payload.ID)
	if err != nil {
		return user, err
	}

	// Commit Transaction
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	user = &models.UserWithRoleResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Roles:     userRoles,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return user, nil
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
