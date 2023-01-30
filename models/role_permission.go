package models

import (
	"time"

	"gorm.io/datatypes"
)

type RolePermission struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	RoleID    uint           `gorm:"column:role_id" json:"roleId"`
	Role      Role           `gorm:"foreignKey:RoleID;constraint:OnUpdate:CASCADE,OnDelete:Cascade" json:"-"`
	Menu      string         `gorm:"column:menu" json:"menu" validate:"required"`
	Action    datatypes.JSON `gorm:"column:action" json:"action" validate:"required"`
	CreatedAt time.Time      `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"column:created_at" json:"updatedAt"`
}

func (RolePermission) TableName() string {
	return "role_permissions"
}
