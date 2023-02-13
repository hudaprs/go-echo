package services

import (
	"echo-rest/helpers"
	"echo-rest/models"
	"echo-rest/queries"
	"echo-rest/structs"

	"gorm.io/gorm"
)

type RoleService struct {
	DB *gorm.DB
}

func (rs *RoleService) Index(pagination helpers.Pagination) (*helpers.Pagination, error) {
	var roles []models.RoleWithPermissionResponse

	query := rs.DB.Scopes(helpers.Paginate(roles, &pagination, rs.DB)).Scopes(queries.RolePermissionPreload()).Find(&roles)
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

	query := rs.DB.Scopes(queries.RolePermissionPreload()).First(&role, id)

	statusCode, err := helpers.ErrorDatabaseNotFound(query.Error)

	return role, statusCode, err
}

func (rs *RoleService) Update(role models.RoleWithPermissionResponse) (models.RoleWithPermissionResponse, error) {
	query := rs.DB.Omit("Permissions.*").Save(&role)

	return role, query.Error
}

func (rs *RoleService) Delete(id uint) (models.RoleWithPermissionResponse, int, error) {
	var role models.RoleWithPermissionResponse

	query := rs.DB.Scopes(queries.RolePermissionPreload()).First(&role, id).Delete(role)

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
