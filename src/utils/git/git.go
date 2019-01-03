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
func GitOps(c *configuration.Config) {

	// Set up git instance
	repo := new(Git)
	repo.Path = c.Workspace

	// Open Repo for use by Contra
	err := GitOpen(repo)
	if err != nil {
		log.Printf("WARNING: Unable to open GIT repository: %v\n", err)
		return
	}

	worktree, err := repo.Repo.Worktree()
	if err != nil {
		log.Printf("WARNING: Failed to fetch GIT working tree: %v\n", err)
		return
	}

	// Grab status and changes
	status, changes, err := GitStatus(*worktree)
	if err != nil {
		log.Printf("WARNING: Failed to check GIT status: %v\n", err)
		return
	}

	// Changes will evaluate to true if something has changed
	if changes {
		// Commit if changes detected
		err := Commit(status, *worktree)
		if err != nil {
			log.Printf("WARNING: Error encountered during GIT commit: %v\n", err)
			return
		}
		log.Println("GIT changes committed.")

		// Attempt to send email. Email failure is logged, but does not interrupt this process.
		gitSendEmail(c, changesOut, changedFiles)

		// push to remote if configured
		if c.GitPush {
			// If private key file is set, init public key auth.
			auth, err := gitSSHAuth(c)
			if err != nil {
				log.Printf(`WARNING: GitPush failed trying to establish GIT authentication: "%s" changes will not be pushed.`, err)
				log.Printf(`INFO: Check configuration settings in %s`, c.ConfigFile)
				log.Printf(`INFO: Set GitAuth = false if not using a custom private key.`)
				log.Printf(`INFO: Set GitPush = false if not using a remote repository.`)
				return
			}
			err = repo.Repo.Push(&git.PushOptions{Auth: auth})
			if err != nil {
				log.Printf("WARNING: Failed to GIT push to remote: %v\n", err)
				return
			}
			log.Println("GIT Push successful.")
		}
	}
}

// gitSendEmail sends git related email notifications
func gitSendEmail(c *configuration.Config, changes, changedFiles []string) {

	// Bail out if email is disabled
	if !c.EmailEnabled {
		// Email is a core feature of this tool, so we are noisy about this being turned off.
		// If other users complain, we can consider silencing this, or adding another parameter
		// to silence it. Requiring two settings in order to quietly skip emailing.
		log.Println("Email notifications are disabled.")
		return
	}

	// Convert slice of changes to a comma separated string
	changesString := strings.Join(diffs, "\n")

	// Send email with changes
	err := utils.SendEmail(c, c.EmailSubject, changesString)
	if err != nil {
		// Log the error, but carry on.
		log.Printf("WARNING: GIT notification email error: %v\n", err)
	}
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
