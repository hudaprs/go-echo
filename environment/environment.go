package environment

import (
	"os"

	"github.com/joho/godotenv"
)

func EnvironmentInit() {
	env := os.Getenv("MODE")
	if env == "" {
		env = "development"
	}

	godotenv.Load(".env." + "env" + ".local")

	if env != "test" {
		godotenv.Load(".env.local")
	}

	godotenv.Load(".env." + env)
	godotenv.Load() // The original .env
}
