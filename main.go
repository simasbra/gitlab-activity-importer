package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

func main() {

	startNow := time.Now()

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
	if len(projectIds) == 0 {
		log.Print("No contributions found for this user. Closing the program.")
		return
	}

	log.Printf("Found contributions in %v projects \n", len(projectIds))

	repo := openOrInitRepo()

	var wg sync.WaitGroup
	commitChannel := make(chan []Commit)

	go func() {
		totalCommits := 0
		for commits := range commitChannel {
			localCommits := createLocalCommit(repo, commits)
			totalCommits += localCommits
		}
		log.Printf("Imported %v commits.\n", totalCommits)

		pushImportedCommits(repo)
	}()

	for _, projectId := range projectIds {
		wg.Add(1)
		go func(projId int) {
			defer wg.Done()

			commits := getProjectCommits(projId, os.Getenv("COMMITER_NAME"))

			commitChannel <- commits

		}(projectId)
	}
	wg.Wait()

	pushImportedCommits(repo)

	fmt.Printf("This operation took: %v \n", time.Since(startNow))
}
