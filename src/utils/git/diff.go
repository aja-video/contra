package utils

import (
	"fmt"
	"github.com/aja-video/contra/src/configuration"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/pmezard/go-difflib/difflib"
	"io"
	"io/ioutil"
	"log"
)

// GitDiff returns a diff of the existing config and the collector output
// TODO:Move this to diffStrings()
func GitDiff(workspace, filename string, output string) (string, error) {
	qualifiedFile := workspace + `/` + filename
	oldFile, err := ioutil.ReadFile(qualifiedFile)
	if err != nil {
		log.Printf("unable to open an existing device config file %s - assuming new device\n", qualifiedFile)
		oldFile = nil
	}
	diff := difflib.UnifiedDiff{
		A: difflib.SplitLines(string(oldFile)),
		B: difflib.SplitLines(output),
		// diff source name
		FromFile: filename,
		// diff destination name
		ToFile:  "Collector Output",
		Context: 3,
	}
	changes, err := difflib.GetUnifiedDiffString(diff)
	return changes, err
}

// gitFileContent returns the contents of a file tracked in git.
// pass in a git commit Tree object to access previous commits
func gitFileContent(tree *object.Tree, path string) (string, error) {
	entry, err := tree.File(path)
	if err != nil {
		return "", err
	}
	return entry.Contents()
}

// diffStrings returns a diff of two strings
func diffStrings(old string, new string, filename string, toFile string) (string, error) {
	diff := difflib.UnifiedDiff{
		A: difflib.SplitLines(old),
		B: difflib.SplitLines(new),
		// diff source name
		FromFile: filename,
		// diff destination name
		ToFile:  toFile,
		Context: 3,
	}
	changes, err := difflib.GetUnifiedDiffString(diff)
	return changes, err
}

// GitRevDiff returns the difference between the current and previous commit for a given file.
// if rev is passed we look for changes on a specific commit, otherwise the last commit
func GitRevDiff(c *configuration.Config, filepath string, rev ...string) (string, error) {
	repo, err := git.PlainOpen(c.Workspace)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	// Get the HEAD reference (latest commit)
	ref, err := repo.Head()
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	// Get the HEAD commit object
	headCommit, err := repo.CommitObject(ref.Hash())
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	// Get the previous commit object (if it exists)
	var previousCommit *object.Commit
	iter := headCommit.Parents()
	defer iter.Close()
	for { // iterate through commits, assign the first parent object unless we're looking for a specific commit
		parent, err := iter.Next()
		if err == io.EOF {
			// no parents found, or we've reached the end of the commits.
			return fmt.Sprintf("No changes found for %s\n", filepath), nil
		}
		if err != nil {
			fmt.Println(err)
			return "", err
		}
		// if a specific commit hash was passed we'll diff against that
		if len(rev) > 0 {
			if parent.Hash.String() == rev[0] {
				previousCommit = parent
				break
			}
		} else {
			// if rev was not specified we go with the last commit
			previousCommit = parent
			break
		}
	}

	previousTree, err := previousCommit.Tree()
	if err != nil {
		return "Error accessing previous commit", err
	}

	headTree, err := headCommit.Tree()
	if err != nil {
		return "Error accessing git repository", err
	}

	currentFile, err := gitFileContent(headTree, filepath)
	if err != nil {
		return "Error reading current file", err
	}

	previousFile, err := gitFileContent(previousTree, filepath)
	if err != nil {
		return "Error reading previous file", err
	}

	return diffStrings(previousFile, currentFile, filepath, filepath)

}
