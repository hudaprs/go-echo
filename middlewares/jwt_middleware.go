package middlewares

import (
	"echo-rest/helpers"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v4"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func JwtConfig() echojwt.Config {
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(helpers.JwtCustomClaims)
		},
		ErrorHandler: func(c echo.Context, err error) error {
			return helpers.ErrorDynamic(http.StatusUnauthorized, err.Error())
		},
		SigningKey: []byte(os.Getenv("JWT_SECRET")),
	}

	return config
}
