package models

import (
	"echo-rest/database"
	"echo-rest/helpers"
	"time"
)

type Role struct {
	ID             uint             `gorm:"primaryKey" json:"id"`
	Name           string           `gorm:"column:name" json:"name"`
	RolePermission []RolePermission `gorm:"foreignKey:RoleID" json:"permissions"`
	CreatedAt      time.Time        `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt      time.Time        `gorm:"column:updated_at" json:"updatedAt"`
}

type RoleCreateForm struct {
	Name string `json:"name" validate:"required"`
}

type RoleUpdateForm struct {
	Name string `json:"name" validate:"required"`
}

type RoleUpdatePermissionForm struct {
	Permissions []RolePermission `json:"permissions" validate:"required,dive,required"`
}

func (Role) TableName() string {
	return "roles"
}

func (r *Role) MapResponse() []RolePermission {
	return r.RolePermission
}

func (r Role) Index(pagination helpers.Pagination) (*helpers.Pagination, error) {
	db := database.Connect()
	var roles []Role

	query := db.Scopes(helpers.Paginate(roles, &pagination, db)).Preload("RolePermission").Find(&roles)
	pagination.Rows = roles

	return &pagination, query.Error
}

func (r Role) Store(payload RoleCreateForm) (Role, error) {
	db := database.Connect()
	role := Role{
		Name:           payload.Name,
		RolePermission: []RolePermission{},
	}

	query := db.Create(&role)

	return role, query.Error
}

func (r Role) Show(id uint) (Role, int, error) {
	db := database.Connect()

	query := db.Preload("RolePermission").First(&r, id)
	r.MapResponse()

	statusCode, err := helpers.ErrorDatabaseNotFound(query.Error)

	return r, statusCode, err
}

func (Role) Update(role Role) (Role, error) {
	db := database.Connect()

	query := db.Preload("RolePermission").Save(&role)

	return role, query.Error
}

func (Role) UpdatePermission(role Role) (Role, error) {
	db := database.Connect()

	err := db.Where(&Role{ID: role.ID}).Association("RolePermission").Replace(role)

	return role, err
}

func (r Role) Delete(id uint) (Role, int, error) {
	db := database.Connect()

	query := db.First(&r, id).Delete(r)

	statusCode, err := helpers.ErrorDatabaseNotFound(query.Error)

	return r, statusCode, err
}
