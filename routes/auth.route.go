package routes

import (
	"go-echo/controllers"
	"go-echo/middlewares"
)

func AuthRoute(routes RouteGroup, controller controllers.AuthController) {
	auth := routes.V1.Group("/auth")
	auth.POST("/register", controller.Register)
	auth.POST("/login", controller.Login)
	auth.GET("/refresh", controller.Refresh)
	auth.GET("/logout", controller.Logout)
	auth.GET("/me", controller.Me, middlewares.AuthCheck())
	auth.PATCH("/roles/activate/:roleId", controller.ActivateRole, middlewares.AuthCheck())
}
