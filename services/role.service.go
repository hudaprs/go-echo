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
	var roles []models.Role

	query := rs.DB.Scopes(helpers.Paginate(roles, &pagination, rs.DB)).Find(&roles)
	pagination.Rows = roles

	return &pagination, query.Error
}

func (rs *RoleService) Store(payload structs.RoleCreateEditForm) (models.Role, error) {
	role := models.Role{
		Name: payload.Name,
	}

	query := rs.DB.Create(&role)

	return role, query.Error
}

func (rs *RoleService) Show(id uint) (models.Role, int, error) {
	var role models.Role

	query := rs.DB.First(&role, id)

	statusCode, err := helpers.ErrorDatabaseNotFound(query.Error)

	return role, statusCode, err
}

func (rs *RoleService) Update(role models.Role) (models.Role, error) {
	query := rs.DB.Save(&role)

	return role, query.Error
}

func (rs *RoleService) Delete(id uint) (models.Role, int, error) {
	var role models.Role

	query := rs.DB.First(&role, id).Delete(role)

	statusCode, err := helpers.ErrorDatabaseNotFound(query.Error)

	return role, statusCode, err
}
