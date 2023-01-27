package models

import (
	"echo-rest/database"
	"echo-rest/helpers"
	"time"

	"gorm.io/datatypes"
)

type Role struct {
	ID          uint                             `gorm:"primaryKey" json:"id"`
	Name        string                           `gorm:"column:name" json:"name"`
	Permissions datatypes.JSONType[[]Permission] `gorm:"column:permissions" json:"permissions"`
	CreatedAt   time.Time                        `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt   time.Time                        `gorm:"column:updated_at" json:"updatedAt"`
}

type RoleCreateForm struct {
	Name string `json:"name" validate:"required"`
}

type RoleUpdateForm struct {
	Name        string                           `json:"name" validate:"required"`
	Permissions datatypes.JSONType[[]Permission] `json:"permissions" validate:"required"`
}

func (Role) TableName() string {
	return "roles"
}

func (Role) Index(pagination helpers.Pagination) (*helpers.Pagination, error) {
	db := database.Connect()
	var roles []Role

	query := db.Scopes(helpers.Paginate(roles, &pagination, db)).Find(&roles)
	pagination.Rows = roles

	return &pagination, query.Error
}

func (Role) Store(payload RoleCreateForm) (Role, error) {
	db := database.Connect()
	permission := Permission{}
	role := Role{
		Name:        payload.Name,
		Permissions: permission.GeneratePermissions(),
	}

	query := db.Create(&role)

	return role, query.Error
}

func (Role) Show(id uint) (Role, int, error) {
	db := database.Connect()
	var role Role

	query := db.First(&role, id)

	statusCode, err := helpers.ErrorDatabaseNotFound(query.Error)

	return role, statusCode, err
}

func (Role) Update(role Role) (Role, error) {
	db := database.Connect()

	query := db.Save(&role)

	return role, query.Error
}

func (Role) Delete(id uint) (Role, int, error) {
	db := database.Connect()
	var role Role

	query := db.First(&role, id).Delete(role)

	statusCode, err := helpers.ErrorDatabaseNotFound(query.Error)

	return role, statusCode, err
}
