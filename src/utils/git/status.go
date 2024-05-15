package utils

import (
	"github.com/go-git/go-git/v5"
)

// GitStatus reports the current working tree status
func GitStatus(worktree git.Worktree) (git.Status, bool, error) {
	// Assume no changes
	gitChanged := false

	// Get status
	changes, err := worktree.Status()

	if err != nil {
		return nil, false, err
	}

	// Return true if something has changed
	if !changes.IsClean() {
		gitChanged = true
	}

	return changes, gitChanged, err

}
