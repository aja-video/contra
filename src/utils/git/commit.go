package utils

import (
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"log"
	"time"
)

// Commit will add and commit changes
func Commit(path string, status git.Status, worktree git.Worktree) ([]string, []string, error) {
	// Iterate over changed files to determine what is changed
	var changes []string
	var changedFiles []string
	for file, status := range status {
		// Tack on files
		changedFiles = append(changedFiles, file)

		// Tack on diffs.
		diff, err := GitDiff(path, file)
		if err != nil {
			return nil, nil, err
		}
		changes = append(changes, diff)

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
		return nil, nil, err
	}
	log.Printf("Contra Git Commit: %s", commit)
	return changes, changedFiles, err
}
