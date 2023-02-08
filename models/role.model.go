package models

import (
	"time"
)

type Role struct {
	ID           uint         `gorm:"primaryKey"`
	Name         string       `gorm:"column:name"`
	PermissionID uint         `gorm:"-"`
	Permissions  []Permission `gorm:"many2many:role_permissions;foreignKey:ID;joinForeignKey:RoleID;References:ID;joinReferences:PermissionID"`
	UserID       uint         `gorm:"-"`
	Users        []User       `gorm:"many2many:role_users;foreignKey:ID;joinForeignKey:RoleID;References:ID;joinReferences:UserID"`
	CreatedAt    time.Time    `gorm:"column:created_at"`
	UpdatedAt    time.Time    `gorm:"column:updated_at"`
}

type RoleResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type RoleWithPermissionResponse struct {
	ID           uint                            `json:"id"`
	Name         string                          `json:"name"`
	PermissionID uint                            `gorm:"-" json:"-"`
	Permissions  []PermissionWithActionsResponse `gorm:"many2many:role_permissions;foreignKey:ID;joinForeignKey:RoleID;References:ID;joinReferences:PermissionID" json:"permissions"`
	CreatedAt    time.Time                       `json:"createdAt"`
	UpdatedAt    time.Time                       `json:"updatedAt"`
}

func (Role) TableName() string {
	return "roles"
}

func (RoleResponse) TableName() string {
	return "roles"
}

func (RoleWithPermissionResponse) TableName() string {
	return "roles"
}
