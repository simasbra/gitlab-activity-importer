package internal

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func CheckEnvVariables() {
	if os.Getenv("ENV") == "DEVELOPMENT" {
		err := godotenv.Load()

		if err != nil {
			log.Fatalf("Error loading the .env file: %v", err)
		}
	}

	requiredEnvVars := []string{"BASE_URL", "GITLAB_TOKEN", "COMMITER_NAME", "COMMITER_EMAIL", "ORIGIN_REPO_URL", "ORIGIN_TOKEN"}

	for _, envVal := range requiredEnvVars {
		if os.Getenv(envVal) == "" {
			log.Fatalf("Environment variable %s is not set", envVal)
		}
	}
}

func GetHomeDirectory() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Unable to get the user home directory:", err)
	}
	return homeDir
}
