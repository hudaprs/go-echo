package models

import (
	"echo-rest/database"
	"echo-rest/helpers"
	"net/http"
	"time"

	"gorm.io/gorm"
)

type RefreshToken struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	UserID       uint      `gorm:"column:user_id" json:"userId"`
	User         User      `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"user"`
	RefreshToken string    `gorm:"column:refresh_token" json:"refreshToken"`
	CreatedAt    time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt    time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

type RefreshTokenForm struct {
	UserID       uint   `json:"userId"`
	RefreshToken string `json:"refreshToken"`
}

func (RefreshToken) TableName() string {
	return "refresh_tokens"
}

func (rt RefreshToken) Store(payload RefreshTokenForm) (RefreshToken, error) {
	db := database.DatabaseConnection()

	refreshToken := RefreshToken{
		UserID:       payload.UserID,
		RefreshToken: payload.RefreshToken,
	}

	query := db.Create(&refreshToken)

	return refreshToken, query.Error
}

func (RefreshToken) Show(userId uint) (RefreshToken, int, error) {
	db := database.DatabaseConnection()
	var refreshTokenDetail RefreshToken

	query := db.Where(&RefreshToken{UserID: userId}).First(&refreshTokenDetail)

	statusCode, err := helpers.ErrorDatabaseNotFound(query.Error, gorm.ErrRecordNotFound)

	return refreshTokenDetail, statusCode, err
}

func (RefreshToken) ShowByRefreshToken(refreshToken string) (RefreshToken, int, error) {
	db := database.DatabaseConnection()

	var refreshTokenDetail RefreshToken
	query := db.Where(&RefreshToken{RefreshToken: refreshToken}).First(&refreshTokenDetail)

	statusCode, err := helpers.ErrorDatabaseNotFound(query.Error, gorm.ErrRecordNotFound)

	return refreshTokenDetail, statusCode, err
}

func (RefreshToken) Delete(userId uint) error {
	db := database.DatabaseConnection()

	query := db.Where(&RefreshToken{UserID: userId}).Delete(&RefreshToken{})

	return query.Error
}

func (rt RefreshToken) DeleteByRefreshToken(refreshToken string) (int, error) {
	db := database.DatabaseConnection()

	var refreshTokenDetail RefreshToken
	query := db.Where(&RefreshToken{RefreshToken: refreshToken}).First(&refreshTokenDetail)

	if query.Error != nil {
		return http.StatusInternalServerError, query.Error
	}

	_, statusCode, err := rt.Show(refreshTokenDetail.UserID)

	if err != nil {
		return statusCode, err
	}

	queryDelete := db.Delete(&refreshTokenDetail)

	return statusCode, queryDelete.Error
}