package routes

import (
	"echo-rest/controllers"
	"echo-rest/middlewares"
)

func TodoRoute(routes RouteGroup, controller controllers.TodoController) {
	todo := routes.V1.Group("/todos")
	todo.Use(middlewares.AuthCheck())
	todo.GET("", controller.Index, middlewares.RoleCheck("TODO_MANAGEMENT", middlewares.Read))
	todo.GET("/:id", controller.Show, middlewares.RoleCheck("TODO_MANAGEMENT", middlewares.Read))
	todo.POST("", controller.Store, middlewares.RoleCheck("TODO_MANAGEMENT", middlewares.Create))
	todo.PATCH("/:id", controller.Update, middlewares.RoleCheck("TODO_MANAGEMENT", middlewares.Update))
	todo.DELETE("/:id", controller.Delete, middlewares.RoleCheck("TODO_MANAGEMENT", middlewares.Delete))
}
