package main

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"log"
	"os"
)

func openOrInitRepo() *git.Repository {
	repoPath := getHomeDirectory() + "/commits-importer/"
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		if err == git.ErrRepositoryNotExists {
			log.Println("Repository doesn't exist. Initializing a new repository.")
			repo, err = git.PlainInit(repoPath, false)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			log.Fatal("Failed to open or initialize the repository:", err)
		}
	} else {
		log.Println("Opened existing repository.")
	}
	return repo
}

func createLocalCommit(repo *git.Repository, commit []Commit) {
	workTree, err := repo.Worktree()
	if err != nil {
		log.Fatal(err)
	}

	repoPath := getHomeDirectory() + "/commits-importer/"
	filePath := repoPath + "/readme.md"
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		file, err := os.Create(filePath)
		if err != nil {
			log.Fatal(err)
		}
		file.WriteString("Just a readme.")
		file.Close()
	}

	_, err = workTree.Add("readme.md")
	if err != nil {
		log.Fatal(err)
	}

	for index := range commit {
		isDuplicate, _ := checkIfCommitExists(repo, commit[index])
		if isDuplicate {
			log.Println("this commit has been already imported. Skipping.")
			return
		}

		newCommit, err := workTree.Commit(commit[index].ID, &git.CommitOptions{
			Author: &object.Signature{
				Name:  os.Getenv("COMMITER_NAME"),
				Email: os.Getenv("COMMITER_EMAIL"),
				When:  commit[index].AuthoredDate,
			},
			Committer: &object.Signature{
				Name:  os.Getenv("COMMITER_NAME"),
				Email: os.Getenv("COMMITER_EMAIL"),
				When:  commit[index].AuthoredDate,
			},
		})
		if err != nil {
			log.Fatal(err)
		}

		obj, err := repo.CommitObject(newCommit)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Created commit: %s\n", obj.Hash)
	}
}

func checkIfCommitExists(repo *git.Repository, commit Commit) (bool, error) {
	ref, err := repo.Reference("HEAD", true)
	if err != nil {
		return false, err
	}

	iter, err := repo.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		return false, err
	}

	exists := false
	err = iter.ForEach(func(c *object.Commit) error {
		if c.Message == commit.ID {
			exists = true
			return fmt.Errorf("duplicate commit found")
		}
		return nil
	})

	if err != nil {
		return false, err
	}

	return exists, nil
}
