package models

import (
	"gorm.io/datatypes"
)

type RoleUser struct {
	ID      uint                       `gorm:"primaryKey"`
	RoleID  uint                       `gorm:"column:role_id;"`
	Role    Role                       `gorm:"foreignKey:RoleID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	UserID  uint                       `gorm:"column:user_id"`
	User    User                       `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Actions datatypes.JSONType[Action] `gorm:"column:actions"`
}

type RoleUserResponse struct {
	ID     uint `json:"id"`
	RoleID uint `json:"roleId"`
	UserID uint `json:"userId"`
}

func (RoleUser) TableName() string {
	return "role_users"
}

func (RoleUserResponse) TableName() string {
	return "role_users"
}
