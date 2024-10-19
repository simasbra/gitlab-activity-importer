package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/furmanp/gitlab-activity-importer/internal"
	"github.com/furmanp/gitlab-activity-importer/internal/services"
)

func main() {
	internal.CheckEnvVariables()

	gitlabUser := services.GetGitlabUser()

	var result map[string]interface{}
	err := json.Unmarshal([]byte(gitlabUser), &result)

	if err != nil {
		log.Fatalf("Error during parsing GitLab user: %v", err)
	}

	gitLabUserId := result["id"].(float64)

	var projectIds []int
	projectIds, err = services.GetUsersProjectsIds(int(gitLabUserId))

	if err != nil {
		log.Fatalf("Error during getting users projects: %v", err)
	}
	if len(projectIds) == 0 {
		log.Print("No contributions found for this user. Closing the program.")
		return
	}

	log.Printf("Found contributions in %v projects \n", len(projectIds))

	repo := services.OpenOrInitRepo()

	commitChannel := make(chan []internal.Commit, len(projectIds))

	go func() {
		totalCommits := 0
		for commits := range commitChannel {
			localCommits := services.CreateLocalCommit(repo, commits)
			totalCommits += localCommits
		}
		log.Printf("Imported %v commits.\n", totalCommits)

	}()

	services.FetchAllCommits(projectIds, os.Getenv("COMMITER_NAME"), commitChannel)

	services.PushLocalCommits(repo)
}
