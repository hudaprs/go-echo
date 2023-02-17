package models

import (
	"time"
)

type Role struct {
	ID          uint             `gorm:"primaryKey"`
	Name        string           `gorm:"column:name"`
	Permissions []RolePermission `gorm:"-;foreignKey:RoleID;references:ID"`
	Users       []RoleUser       `gorm:"-;foreignKey:RoleID;references:ID"`
	CreatedAt   time.Time        `gorm:"column:created_at"`
	UpdatedAt   time.Time        `gorm:"column:updated_at"`
}

type RoleResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type RoleDropdownResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type RoleWithPermissionResponse struct {
	ID          uint                     `json:"id"`
	Name        string                   `json:"name"`
	Permissions []RolePermissionResponse `gorm:"foreignKey:RoleID;references:ID" json:"permissions"`
	CreatedAt   time.Time                `json:"createdAt"`
	UpdatedAt   time.Time                `json:"updatedAt"`
}

func (Role) TableName() string {
	return "roles"
}

func (RoleResponse) TableName() string {
	return "roles"
}

func (RoleDropdownResponse) TableName() string {
	return "roles"
}

func (RoleWithPermissionResponse) TableName() string {
	return "roles"
}
