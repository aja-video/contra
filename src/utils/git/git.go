package utils

import (
	"contra/src/configuration"
	"gopkg.in/src-d/go-git.v4"
	"log"
)

// Git holds git repo data
type Git struct {
	Repo   *git.Repository
	Path   string
	Remote bool
	url    string
}

// GitOps does stuff with git
func GitOps(c *configuration.Config) error {
	// Set up git instance
	repo := new(Git)
	repo.Path = c.Workspace
	repo.Remote = c.GitPush
	// repo.url = c.GitURL TODO: Determine if this should be configurable

	// Open Repo for use by Contra
	err := GitOpen(repo)
	if err != nil {
		panic(err)
	}

	// Grab status
	status, changes, err := GitStatus(repo)

	if status {
		//GitCommit(repo, changes)
		log.Println(changes)
		log.Println("Commit Placeholder")
	}

	// TODO: Make this useful
	if err != nil {
		log.Printf("Unable to get status information from repository %s:\n %s", repo.Path, err)
	} else {
		log.Printf("Current git status: %s", status)
	}
	return err
}
