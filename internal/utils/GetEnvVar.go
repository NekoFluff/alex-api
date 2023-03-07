package utils

import (
	"fmt"
	"log"
	"os"
	"runtime/debug"

	"github.com/joho/godotenv"
)

func init() {
	// Load the .env file in the current directory
	godotenv.Load()
}

func GetEnvVar(name string) string {
	envVar := os.Getenv(name)
	if envVar == "" {
		debug.PrintStack()
		log.Fatal(fmt.Sprintf("$%v must be set", name))
	}
	return envVar
}
