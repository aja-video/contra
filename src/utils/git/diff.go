package utils

import (
	"github.com/pmezard/go-difflib/difflib"
	"io/ioutil"
	"log"
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
	cmd := exec.Command("git", "diff", "-U4", "HEAD", filename)
	// Switch to the workspace folder.
	cmd.Dir = path
	// Run!
	diff, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(diff), err
}

// gitDiffNative returns a diff of the existing config and the collector output
func gitDiffNative(path string, filename string, output string) string {
	qualifiedFile := path + `/` + filename
	oldFile, err := ioutil.ReadFile(qualifiedFile)
	if err != nil {
		log.Fatalf("Unable to open file %s\n", qualifiedFile)
	}
	diff := difflib.UnifiedDiff{
		A:        difflib.SplitLines(string(oldFile)),
		B:        difflib.SplitLines(output),
		FromFile: filename,
		ToFile:   "Collector Output",
		Context:  3,
	}

	changes, _ := difflib.GetUnifiedDiffString(diff)

	return changes
}
