package models

import (
	"time"
)

type RefreshToken struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	UserID       uint      `gorm:"column:user_id" json:"userId"`
	User         User      `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"user"`
	RefreshToken string    `gorm:"column:refresh_token" json:"refreshToken"`
	CreatedAt    time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt    time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

type RefreshTokenResponse struct {
	ID           uint      `json:"id"`
	UserID       uint      `json:"userId"`
	User         User      `json:"user"`
	RefreshToken string    `json:"refreshToken"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

func (RefreshToken) TableName() string {
	return "refresh_tokens"
}

func (RefreshTokenResponse) TableName() string {
	return "refresh_tokens"
}
