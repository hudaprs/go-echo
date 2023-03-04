package middlewares

import (
	"go-echo/helpers"
	"go-echo/locales"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v4"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func handleErrorMessage(err string) string {
	switch err {
	case "missing value in request header":
		return locales.LocalesGet("validation.authorizationHeader")
	default:
		return err
	}
}

func JwtConfig() echojwt.Config {
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(helpers.JwtCustomClaims)
		},
		ErrorHandler: func(c echo.Context, err error) error {
			return helpers.ErrorDynamic(http.StatusUnauthorized, handleErrorMessage(err.Error()))
		},
		SigningKey: []byte(os.Getenv("JWT_SECRET")),
	}

	return config
}
