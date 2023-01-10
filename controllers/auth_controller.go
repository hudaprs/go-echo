package controllers

import (
	"echo-rest/helpers"
	"echo-rest/models"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

type AuthController struct{}

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

	return helpers.Ok(c, http.StatusCreated, "You has been registered successfully", createdUser)
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

	// Set jwt claims
	claims := &helpers.JwtCustomClaims{
		ID:    userDetail.ID,
		Name:  userDetail.Name,
		Email: userDetail.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}

	// Generate JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate token with signed key (JWT_KEY)
	signedToken, err := token.SignedString([]byte("secret"))
	if err != nil {
		return helpers.ErrorBadRequest(err.Error())
	}

	return helpers.Ok(c, http.StatusOK, "You have successfully login", models.UserAuthResponse{
		Token: signedToken,
	})
}
