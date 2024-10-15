package main

import (
	"encoding/json"
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
	log.Printf("Found contributions in %v projects \n", len(projectIds))

	repo := openOrInitRepo()
	totalCommits := 0

	for index := range projectIds {
		commits, _ := getProjectCommits(projectIds[index], os.Getenv("COMMITER_NAME"))
		localCommits := createLocalCommit(repo, commits)
		totalCommits += localCommits
	}
	log.Printf("Imported %v commits. \n", totalCommits)

}
