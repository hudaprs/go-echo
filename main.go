package main

import (
	"echo-rest/database"
	"echo-rest/environment"
	"echo-rest/routes"

	"github.com/labstack/echo/v4"
)

func main() {
	// Instance of echo
	e := echo.New()

	// Init Environment
	environment.EnvironmentInit()

	// Init Database
	database.DatabaseInit()

	// Init Routes
	routes.RoutesInit(e)

	e.Logger.Fatal(e.Start(":8000"))
}
