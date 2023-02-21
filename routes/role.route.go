package routes

import (
	"echo-rest/controllers"
	"echo-rest/middlewares"
)

func RoleRoute(routes RouteGroup, controller controllers.RoleController) {
	role := routes.V1.Group("/roles")
	role.Use(middlewares.AuthCheck())
	role.GET("", controller.Index, middlewares.RoleCheck("ROLE_MANAGEMENT", middlewares.Read))
	role.GET("/:id", controller.Show, middlewares.RoleCheck("ROLE_MANAGEMENT", middlewares.Read))
	role.POST("", controller.Store, middlewares.RoleCheck("ROLE_MANAGEMENT", middlewares.Create))
	role.PATCH("/:id", controller.Update, middlewares.RoleCheck("ROLE_MANAGEMENT", middlewares.Update))
	role.DELETE("/:id", controller.Delete, middlewares.RoleCheck("ROLE_MANAGEMENT", middlewares.Delete))
}
