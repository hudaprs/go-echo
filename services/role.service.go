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

func (rs *RoleService) Index(pagination helpers.Pagination) (*helpers.Pagination, error) {
	var roles []models.RoleResponse

	query := rs.DB.Scopes(helpers.Paginate(roles, &pagination, rs.DB)).Preload("Permissions", func(db *gorm.DB) *gorm.DB {
		return db.Select("permissions.*, role_permissions.actions").Joins("left join role_permissions on role_permissions.permission_id = permissions.id")
	}).Find(&roles)
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

func (rs *RoleService) Show(id uint) (models.RoleResponse, int, error) {
	var role models.RoleResponse

	query := rs.DB.First(&role, id)

	statusCode, err := helpers.ErrorDatabaseNotFound(query.Error)

	return role, statusCode, err
}

func (rs *RoleService) Update(role models.RoleResponse) (models.RoleResponse, error) {
	query := rs.DB.Save(&role)

	return role, query.Error
}

func (rs *RoleService) Delete(id uint) (models.RoleResponse, int, error) {
	var role models.RoleResponse

	query := rs.DB.First(&role, id).Delete(role)

	statusCode, err := helpers.ErrorDatabaseNotFound(query.Error)

	return role, statusCode, err
}
