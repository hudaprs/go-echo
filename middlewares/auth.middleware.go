package middlewares

import (
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func AuthCheck() echo.MiddlewareFunc {
	return echojwt.WithConfig(JwtConfig())
}
