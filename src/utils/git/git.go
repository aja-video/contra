package utils

import (
	"contra/src/configuration"
	"gopkg.in/src-d/go-git.v4"
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
	// Determine if we are going to do a git push
	repo.Remote = c.GitPush

	// Open Repo for use by Contra
	err := GitOpen(repo)

	if err != nil {
		return err
	}

	worktree, err := repo.Repo.Worktree()

	if err != nil {
		return err
	}

	// Grab status and changes
	status, changes, err := GitStatus(*worktree)
	// Status will evaluate to true if something has changed
	if changes {
		// Commit if changes detected
		err = Commit(status, *worktree)
		//TODO: Diffs
		// push to remote if configured
		if repo.Remote {
			err = repo.Repo.Push(&git.PushOptions{})
		}
	}
	return err
}
