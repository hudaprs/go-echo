package services

import (
	"echo-rest/helpers"
	"echo-rest/models"

	"gorm.io/gorm"
)

type UserService struct {
	DB *gorm.DB
}

func mapRoles(db *gorm.DB) *gorm.DB {
	return db.Select("roles.*").Joins("left join role_users on role_users.role_id = roles.role_id")
}

func (us *UserService) Index(pagination helpers.Pagination) (*helpers.Pagination, error) {
	var users []models.UserWithRoleResponse

	query := us.DB.Scopes(helpers.Paginate(users, &pagination, us.DB)).Preload("Roles", mapRoles).Find(&users)
	pagination.Rows = users

	return &pagination, query.Error
}

func (us *UserService) Store() {
	//
}

func (us *UserService) Show() {
	//
}

func (us *UserService) Update() {
	//
}

func (us *UserService) Delete() {
	//
}
