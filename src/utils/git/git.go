package utils

import (
	"github.com/aja-video/contra/src/configuration"
	"github.com/aja-video/contra/src/utils"
	"gopkg.in/src-d/go-git.v4"
	gitSsh "gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
	"log"
	"strings"
)

// Git holds git repo data
type Git struct {
	Repo *git.Repository
	Path string
	url  string
}

// GitOps does stuff with git
func GitOps(c *configuration.Config) error {

	// Set up git instance
	repo := new(Git)
	repo.Path = c.Workspace

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
		err := Commit(status, *worktree)
		if err != nil {
			return err
		}
		log.Println("GIT changes committed.")

		// push to remote if configured
		if c.GitPush {
			// If private key file is set, init public key auth.
			auth, err := gitSSHAuth(c)
			if err != nil {
				return err
			}
			err = repo.Repo.Push(&git.PushOptions{Auth: auth})
			if err != nil {
				return err
			}
			log.Println("GIT Push successful.")
		}
	}

	return err
}

// GitSendEmail sends git related email notifications
func GitSendEmail(c *configuration.Config, diffs []string) error {

	// Bail out if email is disabled
	if !c.EmailEnabled {
		log.Println("Email notifications are disabled.")
		return nil
	}

	// Convert slice of changes to a comma separated string
	changesString := strings.Join(diffs, "\n")

	// Send email with changes
	err := utils.SendEmail(c, c.EmailSubject, changesString)

	if err != nil {
		return err
	}

	return nil
}

// gitSSHAuth sets up authentication for git a git remote
func gitSSHAuth(c *configuration.Config) (gitSsh.AuthMethod, error) {
	// If auth is disabled, there is nothing to do here.
	if !c.GitAuth {
		return nil, nil
	}

	auth, err := gitSsh.NewPublicKeysFromFile(c.GitUser, c.GitPrivateKey, "")
	return auth, err

}
