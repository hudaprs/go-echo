package services

import (
	"echo-rest/helpers"
	"echo-rest/models"
	"echo-rest/queries"
	"echo-rest/structs"
	"errors"

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
	var mergedRoleList []models.RoleUserResponse

	// Check if theres any permissions existed before
	// Delete all data, and create an new one
	query := rs.DB.Where(models.RoleUserResponse{UserID: userId}).Delete(&models.RoleUser{})
	if query.Error != nil {
		return []models.RoleUserResponse{}, query.Error
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
		return []models.RoleUserResponse{}, query.Error
	}

	if len(mergedRoleList) > 0 {
		return mergedRoleList, query.Error
	} else {
		return []models.RoleUserResponse{}, query.Error
	}
}

func (rs *RoleService) ActivateRole(roleId uint, userId uint) (models.RoleResponse, int, error) {
	var roleMappingDetail *models.RoleUser
	var roleUserList []models.RoleUser
	var roleDetail models.RoleResponse
	var err error
	var statusCode int

	if query := rs.DB.Where(&models.RoleUser{RoleID: roleId, UserID: userId}).First(&roleMappingDetail); query.Error != nil {
		err = query.Error
	}

	// Check if user have specific role
	if roleMappingDetail != nil {
		// Check if user have roles
		if query := rs.DB.Where(&models.RoleUser{UserID: userId}).Not(models.RoleUser{RoleID: roleId}).Find(&roleUserList); query.Error != nil {
			err = query.Error
		}

		// If user have any roles, make other role status inactive
		if len(roleUserList) > 0 {
			for _, roleUser := range roleUserList {
				roleUser.IsActive = false
				if query := rs.DB.Save(&roleUser); query.Error != nil {
					err = query.Error
				}
			}
		}

		// Update the selected role to be active
		roleMappingDetail.IsActive = true
		if query := rs.DB.Save(&roleMappingDetail); query.Error != nil {
			err = query.Error
		}

		// Assign role detail to be selected role by the user
		if query := rs.DB.First(&roleDetail, roleMappingDetail.RoleID); query.Error != nil {
			err = query.Error
		}
	}

	if err != nil {
		statusCode = 500
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		statusCode = 404
	}

	return roleDetail, statusCode, err
}
