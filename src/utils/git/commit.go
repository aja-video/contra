package utils

import (
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"log"
	"time"
)

// Commit will add and commit changes
func Commit(status git.Status, worktree git.Worktree) error {
	log.Println("------------------GIT COMMIT OUTPUT-------------")

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
	return err
}
