package main

import (
	"log"
	"os"
)

func checkEnvVariables() {
	// err := godotenv.Load()

	// if err != nil {
	// 	log.Fatal("Error loading the .env file")
	// }

	requiredEnvVars := []string{"BASE_URL", "GITLAB_TOKEN", "COMMITER_NAME", "COMMITER_EMAIL", "ORIGIN_REPO_URL", "ORIGIN_TOKEN"}

	for _, envVal := range requiredEnvVars {
		if os.Getenv(envVal) == "" {
			log.Fatalf("Environment variable %s is not set", envVal)
		}
	}
}

func getHomeDirectory() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Unable to get the user home directory:", err)
	}
	return homeDir
}
