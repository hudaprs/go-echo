package routes

import (
	"go-echo/controllers"
	"go-echo/middlewares"
)

func UserRoute(routes RouteGroup, controller controllers.UserController) {
	user := routes.V1.Group("/users")
	user.Use(middlewares.AuthCheck())
	user.GET("", controller.Index, middlewares.RoleCheck("USER_MANAGEMENT", middlewares.Read))
	user.GET("/:id", controller.Show, middlewares.RoleCheck("USER_MANAGEMENT", middlewares.Read))
	user.POST("", controller.Store, middlewares.RoleCheck("USER_MANAGEMENT", middlewares.Create))
	user.PATCH("/:id", controller.Update, middlewares.RoleCheck("USER_MANAGEMENT", middlewares.Update))
	user.DELETE("/:id", controller.Delete, middlewares.RoleCheck("USER_MANAGEMENT", middlewares.Delete))
	user.GET("/roles", controller.RoleDropdown, middlewares.RoleCheck("USER_MANAGEMENT", middlewares.Read))
}
