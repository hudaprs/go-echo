package models

import (
	"time"
)

type RoleUser struct {
	ID        uint      `gorm:"primaryKey"`
	RoleID    uint      `gorm:"column:role_id"`
	Role      Role      `gorm:"foreignKey:RoleID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	UserID    uint      `gorm:"column:user_id"`
	User      User      `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	IsActive  bool      `gorm:"column:is_active;default:false"`
	CreatedAt time.Time `gorm:"created_at"`
	UpdatedAt time.Time `gorm:"updated_at"`
}

type RoleUserResponse struct {
	ID        uint      `json:"id"`
	RoleID    uint      `json:"-"`
	UserID    uint      `json:"-"`
	Name      string    `gorm:"<-:false" json:"name"`
	IsActive  bool      `json:"isActive"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type RoleUserWithPermissionResponse struct {
	ID          uint                     `json:"id"`
	RoleID      uint                     `json:"-"`
	UserID      uint                     `json:"-"`
	Name        string                   `gorm:"<-:false" json:"name"`
	Permissions []RolePermissionResponse `gorm:"<-:false;foreignKey:RoleID;references:RoleID" json:"permissions"`
	IsActive    bool                     `json:"isActive"`
	CreatedAt   time.Time                `json:"createdAt"`
	UpdatedAt   time.Time                `json:"updatedAt"`
}

func (RoleUser) TableName() string {
	return "role_users"
}

func (RoleUserResponse) TableName() string {
	return "role_users"
}

func (RoleUserWithPermissionResponse) TableName() string {
	return "role_users"
}
