package services

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/furmanp/gitlab-activity-importer/internal"
)

func GetGitlabUser() string {
	url := os.Getenv("BASE_URL")

	client := &http.Client{}
	req, _ := http.NewRequest("GET", fmt.Sprintf("%v/api/v4/user", url), nil)
	req.Header.Set("PRIVATE-TOKEN", os.Getenv("GITLAB_TOKEN"))

	res, err := client.Do(req)

	if err != nil {
		log.Print("something went wrong with your request", err)
	}

	if res.StatusCode == http.StatusOK {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			log.Fatal("something went wrong")
		}
		json := string(body)

		return json
	}

	return "User not found"
}

func GetUsersProjectsIds(userId int) ([]int, error) {
	url := os.Getenv("BASE_URL")

	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%v/api/v4/users/%v/contributed_projects", url, userId), nil)
	if err != nil {
		log.Fatalf("Error creating the request: %v", err)
	}

	req.Header.Set("PRIVATE-TOKEN", os.Getenv("GITLAB_TOKEN"))
	res, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error making the request: %v", err)
	}

	if res.StatusCode != http.StatusOK {
		log.Fatalf("Request failed with status code: %v", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("Error reading the response body: %v", err)
	}

	res.Body.Close()

	var result []map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}

	if len(result) == 0 {
		log.Fatalf("No contributed projects found")
	}

	var projectIds []int
	for index := range result {
		id := result[index]["id"].(float64)
		projectIds = append(projectIds, int(id))
	}

	return projectIds, nil
}

func GetProjectCommits(projectId int, userName string) []internal.Commit {
	url := os.Getenv("BASE_URL")
	token := os.Getenv("GITLAB_TOKEN")

	var allCommits []internal.Commit
	client := &http.Client{}
	page := 1

	since := time.Now().UTC().Add(-24 * time.Hour).Format(time.RFC3339)

	for {
		req, err := http.NewRequest("GET", fmt.Sprintf("%v/api/v4/projects/%v/repository/commits?author=%v&per_page=100&page=%d&since=%v", url, projectId, userName, page, since), nil)
		if err != nil {
			log.Fatalf("Error fetching the commits: %v", err)
		}

		req.Header.Set("PRIVATE-TOKEN", token)
		res, err := client.Do(req)
		if err != nil {
			log.Fatalf("Error making the request: %v", err)
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			log.Fatalf("Request failed with status code: %v", res.StatusCode)
		}

		body, err := io.ReadAll(res.Body)
		if err != nil {
			log.Fatalf("Error reading the response body: %v", err)
		}

		var commits []internal.Commit
		err = json.Unmarshal(body, &commits)
		if err != nil {
			log.Fatalf("Error parsing JSON: %v", err)
		}

		if len(commits) == 0 {
			break
		}

		allCommits = append(allCommits, commits...)

		page++
	}

	if len(allCommits) == 0 {
		log.Printf("Found no commits in project no.:%v \n", projectId)
		return nil
	}

	log.Printf("Found total of %v commits in project no.:%v \n", len(allCommits), projectId)

	return allCommits
}

func FetchAllCommits(projectIds []int, commiterName string, commitChannel chan []internal.Commit) {
	var wg sync.WaitGroup

	for _, projectId := range projectIds {
		wg.Add(1)

		go func(projId int) {
			defer wg.Done()

			commits := GetProjectCommits(projId, commiterName)
			if len(commits) > 0 {
				commitChannel <- commits
			}

		}(projectId)
	}

	wg.Wait()
	close(commitChannel)

}
