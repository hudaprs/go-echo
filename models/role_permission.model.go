package models

import (
	"gorm.io/datatypes"
)

type Action struct {
	Create bool `json:"create"`
	Read   bool `json:"read"`
	Update bool `json:"update"`
	Delete bool `json:"delete"`
}

type RolePermission struct {
	ID           uint                       `gorm:"primaryKey"`
	RoleID       uint                       `gorm:"column:role_id;"`
	Role         Role                       `gorm:"foreignKey:RoleID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	PermissionID uint                       `gorm:"column:permission_id"`
	Permission   Permission                 `gorm:"foreignKey:PermissionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Actions      datatypes.JSONType[Action] `gorm:"column:actions"`
}

type RolePermissionResponse struct {
	ID           uint                       `json:"id"`
	RoleID       uint                       `json:"roleId"`
	PermissionID uint                       `json:"permissionId"`
	Actions      datatypes.JSONType[Action] `json:"actions"`
}

func (RolePermission) TableName() string {
	return "role_permissions"
}

func (RolePermissionResponse) TableName() string {
	return "role_permissions"
}
