package utils

import (
	"github.com/pmezard/go-difflib/difflib"
	"io/ioutil"
	"log"
)

// GitDiff returns a diff of the existing config and the collector output
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
