package models

import (
	"time"
)

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"column:name" json:"name"`
	Email     string    `gorm:"column:email;index:index_email,unique" json:"email"`
	Password  string    `gorm:"column:password" json:"-"`
	RoleID    uint      `gorm:"-"`
	Roles     []Role    `gorm:"many2many:role_users;foreignKey:ID;joinForeignKey:UserID;References:ID;joinReferences:RoleID"`
	CreatedAt time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

type UserResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type UserWithRoleResponse struct {
	ID        uint           `json:"id"`
	Name      string         `json:"name"`
	Email     string         `json:"email"`
	Password  string         `json:"-"`
	RoleID    uint           `gorm:"-" json:"-"`
	Roles     []RoleResponse `gorm:"many2many:role_users;foreignKey:ID;joinForeignKey:UserID;References:ID;joinReferences:RoleID" json:"roles"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
}

type UserWithRolePermissionResponse struct {
	ID        uint                         `json:"id"`
	Name      string                       `json:"name"`
	Email     string                       `json:"email"`
	Password  string                       `json:"-"`
	RoleID    uint                         `gorm:"-" json:"-"`
	Roles     []RoleWithPermissionResponse `gorm:"many2many:role_users;foreignKey:ID;joinForeignKey:UserID;References:ID;joinReferences:RoleID" json:"roles"`
	CreatedAt time.Time                    `json:"createdAt"`
	UpdatedAt time.Time                    `json:"updatedAt"`
}

func (User) TableName() string {
	return "users"
}

func (UserResponse) TableName() string {
	return "users"
}

func (UserWithRoleResponse) TableName() string {
	return "users"
}

func (UserWithRolePermissionResponse) TableName() string {
	return "users"
}
