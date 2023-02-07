package models

import (
	"time"
)

type Todo struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"column:user_id"`
	User      User      `gorm:"foreignKey:UserID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT"`
	Title     string    `gorm:"column:title"`
	Completed bool      `gorm:"column:completed"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

type TodoResponse struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"userId"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type TodoUserResponse struct {
	ID        uint         `json:"id"`
	UserID    uint         `json:"userId"`
	User      UserResponse `json:"user"`
	Title     string       `json:"title"`
	Completed bool         `json:"completed"`
	CreatedAt time.Time    `json:"createdAt"`
	UpdatedAt time.Time    `json:"updatedAt"`
}

func (Todo) TableName() string {
	return "todos"
}

func (TodoResponse) TableName() string {
	return "todos"
}

func (TodoUserResponse) TableName() string {
	return "todos"
}
