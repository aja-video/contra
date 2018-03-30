package utils

import (
	"os/exec"
)

// GitDiff returns a diff.
func GitDiff(path, filename string) (string, error) {
	// Explicitly reminding ourselves that we're using exec.
	return gitDiffExec(path, filename)
}

// gitDiffExec uses os/exec to pull a git diff.
func gitDiffExec(path, filename string) (string, error) {
	// Prep the command.
	cmd := exec.Command("git", "diff", "-U4", filename)
	// Switch to the workspace folder.
	cmd.Dir = path
	// Run!
	diff, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(diff), err
}

func gitDiffNative() {
	// Someday...
}
