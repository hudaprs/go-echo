package controllers

import (
	"echo-rest/helpers"
	"echo-rest/models"
	"echo-rest/services"
	"echo-rest/structs"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

type AuthController struct {
	AuthService         services.AuthService
	RefreshTokenService services.RefreshTokenService
}

// @description Generate token (common token / refresh token)
// @scope 		Private
// @param 		user models.UserResponse, bool rememberMe
// @return		string (common token), string (refreshToken), error
func generateToken(user models.UserResponse) (string, string, error) {
	claims := helpers.JwtCustomClaims{
		ID:    user.ID,
		Email: user.Email,
	}
	refreshClaims := claims

	// Make token expire to 1 hour
	claims.RegisteredClaims = jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 1)),
	}

	// Generate JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)

	// Generate token with signed key (JWT_SECRET)
	signedToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", "", err
	}

	// Generate refresh token with signed key (JWT_SECRET_REFRESH)
	signedRefreshToken, err := refreshToken.SignedString([]byte(os.Getenv("JWT_SECRET_REFRESH")))
	if err != nil {
		return "", "", err
	}

	return signedToken, signedRefreshToken, err
}

// @description Get refresh token from custom headers (X-Auth-Refresh-Token)
// @return		string (refresh token) []helpers.ValidationResponse
func getRefreshToken(c echo.Context) (string, []helpers.ValidationResponse) {
	// Headers
	refreshTokenHeader := c.Request().Header["X-Auth-Refresh-Token"]

	// Check if no refresh token header
	if len(refreshTokenHeader) == 0 {
		var validationResponse []helpers.ValidationResponse

		validationResponse = append(validationResponse, helpers.ValidationResponse{
			Field:   "X-Auth-Refresh-Token",
			Message: "X-Auth-Refresh-Token header is required",
		})

		return "", validationResponse
	}

	return refreshTokenHeader[0], nil
}

// @description Register
// @param 		echo.Context
// @return		error
func (ac AuthController) Register(c echo.Context) error {
	form := new(structs.UserStoreForm)

	if err := c.Bind(form); err != nil {
		return helpers.ErrorBadRequest(err.Error())
	}

	if err := c.Validate(form); err != nil {
		return err
	}

	_, userDetailStatusCode, _ := ac.AuthService.CheckEmail(form.Email)

	if userDetailStatusCode == 200 {
		return helpers.ErrorBadRequest("Email already exists")
	}

	createdUser, err := ac.AuthService.Store(*form)

	if err != nil {
		return helpers.ErrorBadRequest(err.Error())
	}

	return helpers.Ok(http.StatusCreated, "You has been registered successfully", createdUser)
}

// @description Login
// @param 		echo.Context
// @return		error
func (ac AuthController) Login(c echo.Context) error {
	form := new(structs.UserLoginForm)

	if err := c.Bind(form); err != nil {
		return helpers.ErrorBadRequest(err.Error())
	}

	if err := c.Validate(form); err != nil {
		return err
	}

	userDetail, userDetailStatusCode, _ := ac.AuthService.CheckEmail(form.Email)

	// Check if user not exists
	if userDetailStatusCode == 404 {
		return helpers.ErrorBadRequest("Invalid credentials")
	}

	// Check if user exists
	if userDetailStatusCode == 200 {
		// Then compare the password
		isPasswordCorrect := helpers.PasswordCompare(form.Password, userDetail.Password)

		// If password not correct, throw the error
		if !isPasswordCorrect {
			return helpers.ErrorBadRequest("Invalid credentials")
		}
	}

	// Generate token
	token, refreshToken, err := generateToken(userDetail)

	if err != nil {
		return helpers.ErrorBadRequest(err.Error())
	}

	// Find refresh token
	_, refreshTokenDetailStatusCode, refreshTokenDetailErr := ac.RefreshTokenService.Show(userDetail.ID)
	if refreshTokenDetailErr != nil && refreshTokenDetailStatusCode == http.StatusInternalServerError {
		return helpers.ErrorServer(refreshTokenDetailErr.Error())
	}

	// Check if refresh token exists
	if refreshTokenDetailStatusCode == http.StatusOK {
		// Remove refresh token
		refreshTokenDeleteErr := ac.RefreshTokenService.Delete(userDetail.ID)
		if refreshTokenDeleteErr != nil {
			return helpers.ErrorServer(refreshTokenDeleteErr.Error())
		}
	}

	// Insert refresh token to database
	// In other word, create new refresh token every user login again
	_, refreshTokenInsertErr := ac.RefreshTokenService.Store(structs.RefreshTokenForm{
		UserID:       userDetail.ID,
		RefreshToken: refreshToken,
	})

	if refreshTokenInsertErr != nil {
		return helpers.ErrorServer(refreshTokenInsertErr.Error())
	}

	return helpers.Ok(http.StatusOK, "You have successfully login", structs.UserLoginResponse{
		Token:        token,
		RefreshToken: refreshToken,
	})
}

// @description Get authenticated user
// @param 		echo.Context
// @return		error
func (ac AuthController) Me(c echo.Context) error {
	authenticatedUser := helpers.JwtGetClaims(c)

	// Find User
	userDetail, statusCode, err := ac.AuthService.Show(authenticatedUser.ID)

	if err != nil && statusCode >= 400 {
		return helpers.ErrorDynamic(statusCode, err.Error())
	}

	return helpers.Ok(http.StatusOK, "Hi!", userDetail)
}

// @description Refresh token
// @param 		echo.Context
// @return		error
func (ac AuthController) Refresh(c echo.Context) error {
	// Get refresh token from header
	refreshTokenHeaderString, refreshTokenHeaderErr := getRefreshToken(c)
	if refreshTokenHeaderErr != nil {
		return helpers.ErrorValidation(refreshTokenHeaderErr)
	}

	// Find refresh token from database
	refreshTokenDetail, statusCode, err := ac.RefreshTokenService.ShowByRefreshToken(refreshTokenHeaderString)
	if err != nil && statusCode >= 400 {
		return helpers.ErrorDynamic(statusCode, err.Error())
	}

	// Find user model
	userDetail, statusCode, err := ac.AuthService.Show(refreshTokenDetail.UserID)
	if err != nil && statusCode >= 400 {
		return helpers.ErrorDynamic(statusCode, err.Error())
	}

	// Generate new token
	token, newRefreshToken, err := generateToken(userDetail)
	if err != nil {
		return helpers.ErrorBadRequest(err.Error())
	}

	return helpers.Ok(http.StatusOK, "Token refreshed", structs.UserLoginResponse{
		Token:        token,
		RefreshToken: newRefreshToken,
	})
}

// @description Logout
// @param 		echo.Context
// @return		error
func (ac AuthController) Logout(c echo.Context) error {
	// Get refresh token from header
	refreshTokenHeaderString, refreshTokenHeaderErr := getRefreshToken(c)
	if refreshTokenHeaderErr != nil {
		return helpers.ErrorValidation(refreshTokenHeaderErr)
	}

	// Remove refresh token
	statusCode, refreshTokenDeleteErr := ac.RefreshTokenService.DeleteByRefreshToken(refreshTokenHeaderString)
	if refreshTokenDeleteErr != nil {
		return helpers.ErrorDynamic(statusCode, refreshTokenDeleteErr.Error())
	}

	return helpers.Ok(http.StatusOK, "You have successfully logout", nil)
}
