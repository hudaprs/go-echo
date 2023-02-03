package models

import (
	"time"

	"gorm.io/datatypes"
)

type Permission struct {
	ID        uint      `gorm:"primaryKey"`
	Code      string    `gorm:"column:code"`
	RoleID    uint      `gorm:"-"`
	Roles     []Role    `gorm:"many2many:role_permissions;foreignKey:ID;joinForeignKey:PermissionID;References:ID;joinReferences:RoleID"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

type PermissionResponse struct {
	ID        uint   `json:"id"`
	Code      string `json:"code"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type PermissionWithActionsResponse struct {
	ID        uint                       `json:"id"`
	Code      string                     `json:"code"`
	CreatedAt string                     `json:"createdAt"`
	Actions   datatypes.JSONType[Action] `json:"actions"`
	UpdatedAt string                     `json:"updatedAt"`
}

func (Permission) TableName() string {
	return "permissions"
}

func (PermissionResponse) TableName() string {
	return "permissions"
}

func (PermissionWithActionsResponse) TableName() string {
	return "permissions"
}
