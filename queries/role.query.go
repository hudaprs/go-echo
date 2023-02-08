package queries

import "gorm.io/gorm"

// Get roles
// This is used for roles and users tables
func RolesMap(db *gorm.DB) *gorm.DB {
	return db.Select("roles.*").Joins("left join role_users on role_users.role_id = roles.id")
}
