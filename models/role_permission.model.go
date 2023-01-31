package models

type RolePermission struct {
	RoleID       uint `gorm:"primaryKey;column:role_id"`
	PermissionID uint `gorm:"primaryKey;column:permission_id"`
}

func (RolePermission) TableName() string {
	return "role_permissions"
}
