package utils

import (
	"gopkg.in/src-d/go-git.v4"
	"log"
)

//GitOpen opens or initializes the repository
func GitOpen(g *Git) error {
	var err error
	// Open existing local repo
	g.Repo, err = git.PlainOpen(g.Path)
	// If repo does not exist try to create it
	if err == git.ErrRepositoryNotExists {
		log.Printf("Creating new GIT repository at %s\n", g.Path)
		g.Repo, err = git.PlainInit(
			g.Path, // workspace path
			false,  // isBare = false
		)
		if err != nil {
			// Pass the error up if we cannot create a repo
			return err
		}
	}
	// No error reported
	return nil
}
