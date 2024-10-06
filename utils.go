package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func checkEnvVariables() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading the .env file")
	}

	requiredEnvVars := []string{"BASE_URL", "GITLAB_TOKEN", "COMMITER_NAME", "COMMITER_EMAIL"}

	for _, envVal := range requiredEnvVars {
		if os.Getenv(envVal) == "" {
			log.Fatalf("Environment variable %s is not set", envVal)
		}
	}
}
