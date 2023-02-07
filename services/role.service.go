package services

import (
	"echo-rest/helpers"
	"echo-rest/models"
	"echo-rest/structs"

	"gorm.io/gorm"
)

type RoleService struct {
	DB *gorm.DB
}

func mapPermissions(db *gorm.DB) *gorm.DB {
	return db.Select("permissions.*, role_permissions.actions").Joins("left join role_permissions on role_permissions.permission_id = permissions.id")
}

func (rs *RoleService) Index(pagination helpers.Pagination) (*helpers.Pagination, error) {
	var roles []models.RoleWithPermissionResponse

	query := rs.DB.Scopes(helpers.Paginate(roles, &pagination, rs.DB)).Preload("Permissions", mapPermissions).Find(&roles)
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

	query := rs.DB.Preload("Permissions", mapPermissions).First(&role, id)

	statusCode, err := helpers.ErrorDatabaseNotFound(query.Error)

	return role, statusCode, err
}

func (rs *RoleService) Update(role models.RoleWithPermissionResponse) (models.RoleWithPermissionResponse, error) {
	query := rs.DB.Omit("Permissions.*").Save(&role)

	return role, query.Error
}

func (rs *RoleService) Delete(id uint) (models.RoleWithPermissionResponse, int, error) {
	var role models.RoleWithPermissionResponse

	query := rs.DB.Preload("Permissions", mapPermissions).First(&role, id).Delete(role)

	statusCode, err := helpers.ErrorDatabaseNotFound(query.Error)

	return role, statusCode, err
}
