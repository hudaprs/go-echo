package queries

import "gorm.io/gorm"

// Preload role of users
// This is used for roles and users tables (many2many)
func RoleUserPreload() func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Preload("Roles", func(preloadDb *gorm.DB) *gorm.DB {
			return preloadDb.Joins("left join roles on roles.id = role_users.role_id").Select("role_users.*, roles.name, roles.id")
		})
	}
}

// Preload role of users with role permissions
// This is used for roles, permissions and users tables (many2many)
func RoleUserWithPermissionPreload() func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Scopes(RoleUserPreload()).Scopes(RolePermissionPreload("Roles.Permissions"))
	}
}
