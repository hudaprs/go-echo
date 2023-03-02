package services

import (
	"go-echo/database"
	"go-echo/helpers"
	"go-echo/models"
	"go-echo/queries"
	"go-echo/structs"

	"gorm.io/gorm"
)

type RoleService struct {
	DB *gorm.DB
}

func (rs *RoleService) Index(pagination helpers.Pagination) (*helpers.Pagination, error) {
	var roles []models.RoleResponse

	query := rs.DB.Scopes(helpers.Paginate(roles, &pagination, rs.DB)).Find(&roles)
	pagination.Rows = roles

	return &pagination, query.Error
}

func (rs *RoleService) Store(payload structs.RoleCreateEditForm) (models.RoleResponse, error) {
	role := models.RoleResponse{
		Name: payload.Name,
	}

	query := rs.DB.Create(&role)

	return role, query.Error
}

func (rs *RoleService) Show(id uint) (models.RoleWithPermissionResponse, int, error) {
	var role models.RoleWithPermissionResponse

	query := rs.DB.Scopes(queries.RolePermissionPreload("Permissions")).First(&role, id)

	statusCode, err := helpers.ErrorDatabaseNotFound(query.Error)

	return role, statusCode, err
}

func (rs *RoleService) Update(role models.RoleWithPermissionResponse) (models.RoleWithPermissionResponse, error) {
	query := rs.DB.Omit("Permissions.*").Save(&role)

	return role, query.Error
}

func (rs *RoleService) Delete(id uint) (models.RoleResponse, int, error) {
	var role models.RoleResponse

	query := rs.DB.First(&role, id).Delete(role)

	statusCode, err := helpers.ErrorDatabaseNotFound(query.Error)

	return role, statusCode, err
}

func (rs *RoleService) AssignRoles(userId uint, payload structs.RoleAssignUsersForm) ([]models.RoleUserResponse, error) {
	mergedRoleList := []models.RoleUserResponse{}

	// Check if theres any permissions existed before
	// Delete all data, and create an new one
	query := rs.DB.Where(models.RoleUserResponse{UserID: userId}).Delete(&models.RoleUser{})
	if query.Error != nil {
		return mergedRoleList, query.Error
	}

	// Create New Permissions
	for _, roleUserPayload := range payload.Roles {
		mergedRoleList = append(mergedRoleList, models.RoleUserResponse{
			UserID: userId,
			RoleID: roleUserPayload,
		})
	}

	query = rs.DB.Create(&mergedRoleList)
	if query.Error != nil {
		return mergedRoleList, query.Error
	}

	return mergedRoleList, query.Error
}

func (rs *RoleService) ActivateRole(roleId uint, userId uint) (*models.RoleResponse, int, error) {
	var roleMappingDetail *models.RoleUser
	var roleDetail models.RoleResponse
	tx, err := database.BeginTransaction(rs.DB)
	if err != nil {
		return nil, 500, err
	}

	// Find role inside mapping table
	// Find role by specific role id and user id
	if err := rs.DB.Where(&models.RoleUser{RoleID: roleId, UserID: userId}).First(&roleMappingDetail).Error; err != nil {
		statusCode, err := helpers.ErrorDatabaseNotFound(err)
		return nil, statusCode, err
	}

	// Check if user have specific role
	if roleMappingDetail != nil {
		// Update the selected role to be active
		roleMappingDetail.IsActive = true
		if err := tx.Save(&roleMappingDetail).Error; err != nil {
			return nil, 500, err
		}

		// Assign role detail to be selected role by the user
		if err := rs.DB.First(&roleDetail, roleMappingDetail.RoleID).Error; err != nil {
			statusCode, err := helpers.ErrorDatabaseNotFound(err)
			return nil, statusCode, err
		}

		// Find roles by specific user
		// Find roles that user have, but not the selected (the user selected role to be activated)
		var roleUserListExceptTheSelected []models.RoleUser
		if err := rs.DB.Where(&models.RoleUser{UserID: userId}).Not(models.RoleUser{RoleID: roleId}).Find(&roleUserListExceptTheSelected).Error; err != nil {
			return nil, 500, err
		}

		// If user have any roles, make other role status inactive
		if len(roleUserListExceptTheSelected) > 0 {
			for _, roleUser := range roleUserListExceptTheSelected {
				roleUser.IsActive = false
				if err := tx.Save(&roleUser).Error; err != nil {
					return nil, 500, err
				}
			}
		}
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return nil, 500, err
	}

	return &roleDetail, 200, nil
}

func (rs *RoleService) Dropdown() ([]models.RoleDropdownResponse, error) {
	roleDropdownList := []models.RoleDropdownResponse{}

	if err := rs.DB.Order("name asc").Find(&roleDropdownList).Error; err != nil {
		return roleDropdownList, err
	}

	return roleDropdownList, nil
}
