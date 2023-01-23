package controllers

import (
	"echo-rest/helpers"
	"echo-rest/models"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

type AuthController struct{}

// @description Generate token (common token / refresh token)
// @scope 		Private
// @param 		user models.User
// @return		string (common token), string (refreshToken), error
func generateToken(user models.User) (string, string, error) {
	claims := helpers.JwtCustomClaims{
		ID:    user.ID,
		Email: user.Email,
	}
	refreshClaims := claims

	// Expire for 15 seconds
	claims.RegisteredClaims = jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * 15)),
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

// @description Register
// @param 		echo.Context
// @return		error
func (AuthController) Register(c echo.Context) error {
	form := new(models.UserStoreForm)

	if err := c.Bind(form); err != nil {
		return helpers.ErrorBadRequest(err.Error())
	}

	if err := c.Validate(form); err != nil {
		return err
	}

	user := models.User{}
	_, userDetailStatusCode, _ := user.CheckEmail(form.Email)

	if userDetailStatusCode == 200 {
		return helpers.ErrorBadRequest("Email already exists")
	}

	createdUser, err := user.Store(*form)

	if err != nil {
		return helpers.ErrorBadRequest(err.Error())
	}

	return helpers.Ok(http.StatusCreated, "You has been registered successfully", createdUser)
}

// @description Login
// @param 		echo.Context
// @return		error
func (AuthController) Login(c echo.Context) error {
	form := new(models.UserLoginForm)

	if err := c.Bind(form); err != nil {
		return helpers.ErrorBadRequest(err.Error())
	}

	if err := c.Validate(form); err != nil {
		return err
	}

	user := models.User{}
	userDetail, userDetailStatusCode, _ := user.CheckEmail(form.Email)

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

	refreshTokenModel := models.RefreshToken{}

	// Find refresh token
	_, refreshTokenDetailStatusCode, refreshTokenDetailErr := refreshTokenModel.Show(userDetail.ID)
	if refreshTokenDetailErr != nil && refreshTokenDetailStatusCode == http.StatusInternalServerError {
		return helpers.ErrorServer(refreshTokenDetailErr.Error())
	}

	// Check if refresh token exists
	if refreshTokenDetailStatusCode == http.StatusOK {
		// Remove refresh token
		refreshTokenDeleteErr := refreshTokenModel.Delete(userDetail.ID)
		if refreshTokenDeleteErr != nil {
			return helpers.ErrorServer(refreshTokenDeleteErr.Error())
		}
	}

	// Insert refresh token to database
	// In other word, create new refresh token every user login again
	_, refreshTokenInsertErr := refreshTokenModel.Store(models.RefreshTokenForm{
		UserID:       userDetail.ID,
		RefreshToken: refreshToken,
	})

	if refreshTokenInsertErr != nil {
		return helpers.ErrorServer(refreshTokenInsertErr.Error())
	}

	return helpers.Ok(http.StatusOK, "You have successfully login", models.UserLoginResponse{
		Token:        token,
		RefreshToken: refreshToken,
	})
}

// @description Get authenticated user
// @param 		echo.Context
// @return		error
func (AuthController) Me(c echo.Context) error {
	authenticatedUser := helpers.JwtGetClaims(c)

	// Find User
	user := models.User{}
	userDetail, statusCode, err := user.GetDetail(int(authenticatedUser.ID))

	if err != nil && statusCode >= 400 {
		return helpers.ErrorDynamic(statusCode, err.Error())
	}

	return helpers.Ok(http.StatusOK, "Hi!", userDetail)
}

// @description Refresh token
// @param 		echo.Context
// @return		error
func (AuthController) Refresh(c echo.Context) error {
	form := new(models.UserRefreshForm)

	if err := c.Bind(form); err != nil {
		return helpers.ErrorBadRequest(err.Error())
	}

	if err := c.Validate(form); err != nil {
		return err
	}

	return helpers.Ok(http.StatusOK, "Token refreshed", models.UserLoginResponse{})
}

// @description Logout
// @param 		echo.Context
// @return		error
func (AuthController) Logout(c echo.Context) error {
	form := new(models.UserRefreshForm)

	if err := c.Bind(form); err != nil {
		return helpers.ErrorBadRequest(err.Error())
	}

	if err := c.Validate(form); err != nil {
		return err
	}

	// Remove refresh token
	refreshToken := models.RefreshToken{}
	statusCode, refreshTokenDeleteErr := refreshToken.DeleteByRefreshToken(form.RefreshToken)
	if refreshTokenDeleteErr != nil {
		return helpers.ErrorDynamic(statusCode, refreshTokenDeleteErr.Error())
	}

	return helpers.Ok(http.StatusOK, "You have successfully logout", nil)
}
