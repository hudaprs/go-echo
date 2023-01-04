package middlewares

import (
	"echo-rest/helpers"

	"github.com/golang-jwt/jwt/v4"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func JwtConfig() echojwt.Config {
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(helpers.JwtCustomClaims)
		},
		SigningKey: []byte("secret"),
	}

	return config
}
