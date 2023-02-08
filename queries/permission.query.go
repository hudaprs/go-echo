package queries

import "gorm.io/gorm"

// Get permissions with actions
// This is used for roles and permissions tables
func PermissionsMap(db *gorm.DB) *gorm.DB {
	return db.Select("permissions.*, role_permissions.actions").Joins("left join role_permissions on role_permissions.permission_id = permissions.id")
}
