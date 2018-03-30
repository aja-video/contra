package utils

import (
	"contra/src/configuration"
	"contra/src/utils"
	"gopkg.in/src-d/go-git.v4"
	"log"
	"strings"
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
		changesOut, err := Commit(repo.Path, status, *worktree)
		if err != nil {
			return err
		}
		err = gitSendEmail(c, changesOut)
		if err != nil {
			return err
		}
		//TODO: Diffs
		// push to remote if configured
		if repo.Remote {
			err = repo.Repo.Push(&git.PushOptions{})
		}
	}
	return err
}

//gitSendEmail sends git related email notifications
func gitSendEmail(c *configuration.Config, changes []string) error {

	// Bail out if email is disabled
	if !c.EmailEnabled {
		log.Println("Email notifications are disabled.")
		return nil
	}

	// Convert slice of changes to a comma separated string
	changesString := strings.Join(changes, "\n")

	log.Printf("%s changed, sending email\n", changesString)

	// Send email with changes
	err := utils.SendEmail(c, "Contra-Changes", changesString)

	if err != nil {
		return err
	}

	return nil
}
