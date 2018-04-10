package utils

import (
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"log"
	"os"
	"time"
)

//GitOpen opens or initializes the repository
func GitOpen(g *Git) error {
	var err error
	// Open existing local repo
	g.Repo, err = git.PlainOpen(g.Path)
	// If repo does not exist try to create it
	if err == git.ErrRepositoryNotExists {
		gitSetupRepo(g)
	}
	// No error reported
	return nil
}

// gitSetupRepo initializes a new git repository and creates an initial commit
func gitSetupRepo(g *Git) error {
	var err error
	log.Printf("Creating new GIT repository at %s\n", g.Path)
	// Initialize repo
	g.Repo, err = git.PlainInit(
		g.Path, // workspace path
		false,  // isBare = false
	)

	// Pass the error up if we cannot create a repo
	if err != nil {
		return err
	}
	// Create README.md and commit
	worktree, err := g.Repo.Worktree()
	// Pass the error up if we cannot open the worktree
	if err != nil {
		return err
	}

	readme, err := os.Create(g.Path + "/README.md")
	// Pass the error up if we cannot create README.md
	if err != nil {
		return err
	}

	// Write out contents of README.md
	readme.WriteString("Contra Configuration Tracking GIT Repository")

	// Add to git
	worktree.Add("README.md")

	// Commit
	worktree.Commit(
		"Contra GIT Setup Commit",
		&git.CommitOptions{
			Author: &object.Signature{
				Name:  "Contra",
				Email: "Contra@example.com",
				When:  time.Now(),
			},
		})

	// Nothing failed
	return nil
}
