package utils

import (
	"github.com/aja-video/contra/src/configuration"
	"github.com/pmezard/go-difflib/difflib"
	"io/ioutil"
	"log"
)

// GitDiff returns a diff of the existing config and the collector output
func GitDiff(config configuration.Config, filename string, output string) (string, error) {
	qualifiedFile := config.Workspace + `/` + filename
	oldFile, err := ioutil.ReadFile(qualifiedFile)
	if err != nil {
		log.Printf("unable to open device configuration %s - assuming new config\n", qualifiedFile)
		oldFile = nil
	}
	diff := difflib.UnifiedDiff{
		A:        difflib.SplitLines(string(oldFile)),
		B:        difflib.SplitLines(output),
		FromFile: filename,
		ToFile:   "Collector Output",
		Context:  3,
	}
	changes, err := difflib.GetUnifiedDiffString(diff)
	return changes, err
}
