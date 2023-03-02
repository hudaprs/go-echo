package services

import (
	"go-echo/database"
	"go-echo/helpers"
	"go-echo/models"
	"go-echo/structs"
	"net/http"

	"gorm.io/gorm"
)

type RefreshTokenService struct {
	DB *gorm.DB
}

func (rts *RefreshTokenService) Store(payload structs.RefreshTokenForm) (models.RefreshTokenResponse, error) {
	db := database.Connect()

	refreshToken := models.RefreshTokenResponse{
		UserID:       payload.UserID,
		RefreshToken: payload.RefreshToken,
	}

	query := db.Create(&refreshToken)

	return refreshToken, query.Error
}

func (rts *RefreshTokenService) Show(userId uint) (models.RefreshTokenResponse, int, error) {
	db := database.Connect()
	var refreshTokenDetail models.RefreshTokenResponse

	query := db.Where(&models.RefreshTokenResponse{UserID: userId}).First(&refreshTokenDetail)

	statusCode, err := helpers.ErrorDatabaseDynamic(query.Error)

	return refreshTokenDetail, statusCode, err
}

func (rts *RefreshTokenService) ShowByRefreshToken(refreshToken string) (models.RefreshTokenResponse, int, error) {
	db := database.Connect()

	var refreshTokenDetail models.RefreshTokenResponse
	query := db.Where(&models.RefreshTokenResponse{RefreshToken: refreshToken}).First(&refreshTokenDetail)

	statusCode, err := helpers.ErrorDatabaseDynamic(query.Error)

	return refreshTokenDetail, statusCode, err
}

func (rts *RefreshTokenService) Delete(userId uint) error {
	db := database.Connect()

	query := db.Where(&models.RefreshTokenResponse{UserID: userId}).Delete(&models.RefreshTokenResponse{})

	return query.Error
}

func (rt *RefreshTokenService) DeleteByRefreshToken(refreshToken string) (int, error) {
	db := database.Connect()

	var refreshTokenDetail models.RefreshTokenResponse
	query := db.Where(&models.RefreshTokenResponse{RefreshToken: refreshToken}).First(&refreshTokenDetail)

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
