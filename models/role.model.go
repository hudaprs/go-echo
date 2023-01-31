package models

import (
	"time"
)

type Role struct {
	ID           uint         `gorm:"primaryKey"`
	Name         string       `gorm:"column:name"`
	PermissionID uint         `gorm:"-"`
	Permissions  []Permission `gorm:"many2many:role_permissions;foreignKey:ID;joinForeignKey:RoleID;References:ID;joinReferences:PermissionID"`
	CreatedAt    time.Time    `gorm:"column:created_at"`
	UpdatedAt    time.Time    `gorm:"column:updated_at"`
}

type RoleResponse struct {
	ID          uint         `json:"id"`
	Name        string       `json:"name"`
	Permissions []Permission `json:"permissions"`
	CreatedAt   time.Time    `json:"createdAt"`
	UpdatedAt   time.Time    `json:"updatedAt"`
}

func (Role) TableName() string {
	return "roles"
}

func (RoleResponse) TableName() string {
	return "roles"
}
