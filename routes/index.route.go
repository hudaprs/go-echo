package routes

import (
	"go-echo/controllers"
	"go-echo/database"
	"go-echo/services"
	"net/http"

	"github.com/labstack/echo/v4"
)

type RouteGroup struct {
	V1 *echo.Group
}

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

	// Prefix Route
	v1 := e.Group("/api/v1")

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello World")
	})

	// Auth Feature
	authController := controllers.AuthController{
		AuthService:         AuthService,
		RefreshTokenService: RefreshTokenService,
		UserService:         UserService,
		RoleService:         RoleService,
	}
	AuthRoute(RouteGroup{V1: v1}, authController)

	// User Feature
	userController := controllers.UserController{
		UserService: UserService,
		RoleService: RoleService,
	}
	UserRoute(RouteGroup{V1: v1}, userController)

	// Todo Feature
	todoController := controllers.TodoController{TodoService: TodoService}
	TodoRoute(RouteGroup{V1: v1}, todoController)

	// Role Feature
	roleController := controllers.RoleController{RoleService: RoleService}
	RoleRoute(RouteGroup{V1: v1}, roleController)

	// Permission Feature
	permissionController := controllers.PermissionController{PermissionService: PermissionService}
	PermissionRoute(RouteGroup{V1: v1}, permissionController)

}
