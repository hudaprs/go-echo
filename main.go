package main

import (
	"echo-rest/database"
	"echo-rest/environment"
	"echo-rest/helpers"
	"echo-rest/migrations"
	"echo-rest/routes"

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

	// Init Database
	db, _ := database.DatabaseInit()

	// Init Migration
	migrations.InitMigration(db)

	// Init Routes
	routes.RoutesInit(e)

	e.Logger.Fatal(e.Start("127.0.0.1:8000"))
}
