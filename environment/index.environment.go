package environment

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func envPath(mode string) string {
	dir := filepath.Join("", ".env")

	if mode == "" {
		return dir
	} else {
		return dir + "." + mode
	}
}

func EnvironmentInit() {
	err := godotenv.Load(envPath(""))
	if err != nil {
		panic("Please define .env and MODE inside of it")
	}
	envMode := os.Getenv("MODE")

	fmt.Println("Environment:", envMode)

	godotenv.Load(envPath(envMode))
}
