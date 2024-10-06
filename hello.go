package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func main() {

	checkEnvVariables()

	gitlabUser := getGitlabUser()

	var result map[string]interface{}
	err := json.Unmarshal([]byte(gitlabUser), &result)

	if err != nil {
		log.Fatalf("Error during parsing GitLab user: %v", err)
	}

	gitLabUserId := result["id"].(float64)

	var projectIds []int

	projectIds, err = getUsersProjectsIds(int(gitLabUserId))
	if err != nil {
		log.Fatalf("Error during getting users projects: %v", err)
	}

	fmt.Println(projectIds[2])
	getProjectCommits(projectIds[2], os.Getenv("COMMITER_NAME"))
}
