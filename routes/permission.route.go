package routes

import (
	"go-echo/controllers"
	"go-echo/middlewares"
)

func PermissionRoute(routes RouteGroup, controller controllers.PermissionController) {
	permission := routes.V1.Group("/permissions")
	permission.Use(middlewares.AuthCheck())
	permission.GET("", controller.Index, middlewares.RoleCheck("ROLE_MANAGEMENT", middlewares.Update))
	permission.PATCH("/assign/:roleId", controller.AssignPermissions, middlewares.RoleCheck("ROLE_MANAGEMENT", middlewares.Update))
}
