package services

import (
	"go-echo/database"
	"go-echo/helpers"
	"go-echo/locales"
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
	if err := tx.Where("user_id = ?", userId).Where("role_id NOT IN ?", roleIds).Delete(&models.RoleUserResponse{}).Error; err != nil {
		return roleUserResponse, err
	}

	// Check if user have existing role
	var existedUserRoles []models.RoleUser
	if err := tx.Where(models.RoleUser{UserID: userId}).Where("role_id IN ?", roleIds).Find(&existedUserRoles).Error; err != nil {
		return roleUserResponse, err
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
			if err := tx.First(&models.Role{}, roleId).Error; err != nil {
				_, _err := helpers.ErrorDatabaseDynamic(err, helpers.DatabaseDynamicMessage{
					NotFound: locales.LocalesGet("role.validation.notFound"),
				})
				return roleUserResponse, _err
			}

			newUserRole := &models.RoleUserResponse{
				RoleID: roleId,
				UserID: userId,
			}

			if err := tx.Create(newUserRole).Error; err != nil {
				return roleUserResponse, err
			}
		}
	}

	// Look up new assigned role through database
	if err := tx.Select("role_users.*, roles.name, roles.id").Where("user_id = ?", userId).Joins("left join roles ON roles.id = role_users.role_id").Order("created_at desc").Find(&roleUserResponse).Error; err != nil {
		return roleUserResponse, err
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
	assignedUserRoles := []models.RoleUserResponse{}

	// Make default has password for the first time creating an user
	hashedPassword, err := helpers.PasswordHash("password")
	if err != nil {
		return nil, err
	}

	err = database.BeginTransaction(us.DB, func(tx *gorm.DB) error {
		// Assign value
		if payload.ID != nil {
			if err := tx.Find(&user, payload.ID).Error; err != nil {
				_, _err := helpers.ErrorDatabaseDynamic(err, helpers.DatabaseDynamicMessage{
					NotFound: locales.LocalesGet("user.validation.notFound"),
				})
				return _err
			}
		}

		user.Name = payload.Name
		user.Email = payload.Email
		if payload.ID == nil {
			user.Password = hashedPassword
		}

		// Create / Update new user
		if err := tx.Save(&user).Error; err != nil {
			return err
		}

		// Assign payload.ID to user.id
		payload.ID = &user.ID

		// Create or Update user role
		userRoles, err := syncRoleUser(tx, payload.Roles, *payload.ID)
		if err != nil {
			return err
		}

		assignedUserRoles = userRoles

		return nil
	})
	if err != nil {
		return nil, err
	}

	user = &models.UserWithRoleResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Roles:     assignedUserRoles,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return user, nil
}

func (us *UserService) Show(payload structs.UserAttrsFind) (models.UserWithRoleResponse, int, error) {
	var user models.UserWithRoleResponse

	query := us.DB.Scopes(queries.RoleUserPreload()).Where("id = ?", payload.ID).Or("name = ?", payload.Name).Or("email = ?", payload.Email).First(&user)
	statusCode, err := helpers.ErrorDatabaseDynamic(query.Error, helpers.DatabaseDynamicMessage{
		NotFound: locales.LocalesGet("user.validation.notFound"),
	})

	return user, statusCode, err
}

func (us *UserService) Delete(payload structs.UserAttrsFind) (models.UserWithRoleResponse, int, error) {
	var user models.UserWithRoleResponse

	query := us.DB.Scopes(queries.RoleUserPreload()).Where("id = ?", payload.ID).Or("name = ?", payload.Name).Or("email = ?", payload.Email).First(&user).Delete(user)
	statusCode, err := helpers.ErrorDatabaseDynamic(query.Error, helpers.DatabaseDynamicMessage{
		NotFound: locales.LocalesGet("user.validation.notFound"),
	})

	return user, statusCode, err
}

func (us *UserService) CheckEmail(email string) (bool, int, models.UserResponse, error) {
	var user models.UserResponse
	query := us.DB.Where("email = ?", email).First(&user)

	statusCode, err := helpers.ErrorDatabaseDynamic(query.Error, helpers.DatabaseDynamicMessage{
		NotFound: locales.LocalesGet("user.validation.notFound"),
	})

	return statusCode == 200, statusCode, user, err
}
