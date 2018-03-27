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
		return err
	}

	worktree, err := repo.Repo.Worktree()

	if err != nil {
		return err
	}

	// Grab status
	status, changes, err := GitStatus(*worktree)

	// Status will evaluate to true if something has changed
	if changes {
		// Commit if changes detected
		err = Commit(status, *worktree)
		//TODO: Diffs
		//TODO: push is untested
		if repo.Remote {
			err = repo.Repo.Push(&git.PushOptions{})

		} else {
			log.Println("No changes to commit")
		}
	}

	return err
}
