package routes

import (
	"echo-rest/controllers"
	"echo-rest/database"
	"echo-rest/middlewares"
	"echo-rest/services"
	"net/http"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func RoutesInit(e *echo.Echo) {
	// Database Connection
	db := database.Connect()

	// Services
	AuthService := services.AuthService{DB: db}
	UserService := services.UserService{DB: db}
	RefreshTokenService := services.RefreshTokenService{DB: db}
	TodoService := services.TodoService{DB: db}
	RoleService := services.RoleService{DB: db}
	PermissionService := services.PermissionService{DB: db}

	// Middleware
	authMiddleware := echojwt.WithConfig(middlewares.JwtConfig())

	// Prefix Route
	v1 := e.Group("/api/v1")

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello World")
	})

	// Auth Feature
	AuthController := controllers.AuthController{
		AuthService:         AuthService,
		RefreshTokenService: RefreshTokenService,
	}
	auth := v1.Group("/auth")
	auth.POST("/register", AuthController.Register)
	auth.POST("/login", AuthController.Login)
	auth.GET("/refresh", AuthController.Refresh)
	auth.GET("/logout", AuthController.Logout)
	auth.GET("/me", AuthController.Me, authMiddleware)

	// User Feature
	UserController := controllers.UserController{UserService: UserService}
	user := v1.Group("/users")
	user.Use(authMiddleware)
	user.GET("", UserController.Index)
	user.GET("/:id", UserController.Show)
	user.POST("", UserController.Store)
	user.PATCH("/:id", UserController.Update)
	user.DELETE("/:id", UserController.Delete)

	// Todo Feature
	TodoController := controllers.TodoController{TodoService: TodoService}
	todos := v1.Group("/todos")
	todos.Use(authMiddleware)
	todos.GET("", TodoController.Index)
	todos.GET("/:id", TodoController.Show)
	todos.POST("", TodoController.Store)
	todos.PATCH("/:id", TodoController.Update)
	todos.DELETE("/:id", TodoController.Delete)

	// Role Feature
	RoleController := controllers.RoleController{RoleService: RoleService}
	role := v1.Group("/roles")
	role.Use(authMiddleware)
	role.GET("", RoleController.Index)
	role.GET("/:id", RoleController.Show)
	role.POST("", RoleController.Store)
	role.PATCH("/:id", RoleController.Update)
	role.DELETE("/:id", RoleController.Delete)

	// Permission Feature
	PermissionController := controllers.PermissionController{PermissionService: PermissionService}
	permission := v1.Group("/permissions")
	permission.Use(authMiddleware)
	permission.GET("", PermissionController.Index)
	permission.PATCH("/assign/:roleId", PermissionController.AssignPermissions)
}
