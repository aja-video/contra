package utils

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"log"
	"time"
)

// Commit will add and commit changes
func Commit(status git.Status, worktree git.Worktree) error {
	// Iterate over changed files to determine what is changed
	for file, status := range status {

		// TODO: Maybe a cleaner way to do this?
		switch status.Worktree {
		case git.Untracked:
			log.Printf("New Config File %s\n", file)
			worktree.Add(file)
		case git.Modified:
			log.Printf("Modified Config File %s\n", file)
			worktree.Add(file)
		case git.Deleted:
			log.Printf("Deleted Config File %s\n", file)
			worktree.Remove(file)
		default:
			log.Printf("Unhandled git status for file %s\n", file)

		}
	}

	// Do the commit
	commit, err := worktree.Commit(
		"Contra commit.",
		&git.CommitOptions{
			Author: &object.Signature{
				Name:  "Contra",
				Email: "Contra@example.com",
				When:  time.Now(),
			},
		})
	if err != nil {
		return err
	}
	log.Printf("Contra Git Commit: %s", commit)
	return nil
}
