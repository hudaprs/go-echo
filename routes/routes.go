package routes

import (
	"echo-rest/controllers"
	"echo-rest/middlewares"
	"net/http"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func RoutesInit(e *echo.Echo) {
	v1 := e.Group("/api/v1")

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello World")
	})

	// Auth Features
	AuthController := controllers.AuthController{}
	auth := v1.Group("/auth")
	auth.POST("/register", AuthController.Register)
	auth.POST("/login", AuthController.Login)

	// Todo Features
	TodoController := controllers.TodoController{}
	todos := v1.Group("/todos")
	todos.Use(echojwt.WithConfig(middlewares.JwtConfig()))

	// Todo Route List
	todos.GET("", TodoController.Index)
	todos.GET("/:id", TodoController.Show)
	todos.POST("", TodoController.Store)
	todos.PATCH("/:id", TodoController.Update)
	todos.DELETE("/:id", TodoController.Delete)
}
