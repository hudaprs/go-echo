package models

import (
	"time"

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
	CreatedAt    time.Time                  `gorm:"created_at"`
	UpdatedAt    time.Time                  `gorm:"updated_at"`
}

type RolePermissionResponse struct {
	ID           uint                       `json:"id"`
	RoleID       uint                       `json:"-"`
	PermissionID uint                       `json:"-"`
	Code         string                     `gorm:"<-:false" json:"code"`
	Actions      datatypes.JSONType[Action] `json:"actions"`
	CreatedAt    time.Time                  `json:"createdAt"`
	UpdatedAt    time.Time                  `json:"updatedAt"`
}

func (RolePermission) TableName() string {
	return "role_permissions"
}

func (RolePermissionResponse) TableName() string {
	return "role_permissions"
}
