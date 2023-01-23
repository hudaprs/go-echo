package main

import (
	"echo-rest/database"
	"echo-rest/environment"
	"echo-rest/helpers"
	"echo-rest/migrations"
	"echo-rest/routes"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Instance of echo
	e := echo.New()

	// Validator Instance
	e.Validator = &helpers.CustomValidator{Validator: validator.New()}

	// CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"*"},
	}))

	// Init Environment
	environment.EnvironmentInit()

	// Check if no ENV
	if os.Getenv("DB_HOST") == "" {
		panic("DB_HOST need to be declared in env")
	}
	if os.Getenv("DB_USERNAME") == "" {
		panic("DB_USERNAME need to be declared in env")
	}
	if os.Getenv("DB_PASSWORD") == "" {
		panic("DB_PASSWORD need to be declared in env")
	}
	if os.Getenv("DB_PORT") == "" {
		panic("DB_PORT need to be declared in env")
	}
	if os.Getenv("DB_NAME") == "" {
		panic("DB_NAME need to be declared in env")
	}
	if os.Getenv("JWT_SECRET") == "" {
		panic("JWT_SECRET need to be declared in env")
	}
	if os.Getenv("JWT_SECRET_REFRESH") == "" {
		panic("JWT_SECRET_REFRESH need to be declared in env")
	}

	// Init Database
	db, _ := database.DatabaseInit()

	// Init Migration
	migrations.InitMigration(db)

	// Init Routes
	routes.RoutesInit(e)

	e.Logger.Fatal(e.Start("127.0.0.1:8000"))
}
