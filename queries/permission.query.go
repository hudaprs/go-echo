package queries

import "gorm.io/gorm"

// Preload permission of roles
// This is used for roles and users tables (many2many)
func RolePermissionPreload() func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Preload("Permissions", func(preloadDb *gorm.DB) *gorm.DB {
			return preloadDb.Joins("left join permissions on permissions.id = role_permissions.permission_id").Select("role_permissions.*, permissions.code, permissions.id")
		})
	}
}
