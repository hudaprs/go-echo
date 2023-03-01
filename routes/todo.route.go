package routes

import (
	"go-echo/controllers"
	"go-echo/middlewares"
)

func TodoRoute(routes RouteGroup, controller controllers.TodoController) {
	todo := routes.V1.Group("/todos")
	todo.Use(middlewares.AuthCheck())
	todo.GET("", controller.Index, middlewares.RoleCheck("TODO", middlewares.Read))
	todo.GET("/:id", controller.Show, middlewares.RoleCheck("TODO", middlewares.Read))
	todo.POST("", controller.Store, middlewares.RoleCheck("TODO", middlewares.Create))
	todo.PATCH("/:id", controller.Update, middlewares.RoleCheck("TODO", middlewares.Update))
	todo.DELETE("/:id", controller.Delete, middlewares.RoleCheck("TODO", middlewares.Delete))
}
