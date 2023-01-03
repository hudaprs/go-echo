package routes

import (
	"echo-rest/controllers"
	"net/http"

	"github.com/labstack/echo/v4"
)

func RoutesInit(e *echo.Echo) {
	v1 := e.Group("/api/v1")

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello World")
	})

	// Todo Features
	TodoController := controllers.TodoController{}
	todos := v1.Group("/todos")
	todos.GET("", TodoController.Index)
	todos.GET("/:id", TodoController.Show)
	todos.POST("", TodoController.Store)
	todos.PATCH("/:id", TodoController.Update)
	todos.DELETE("/:id", TodoController.Delete)
}
