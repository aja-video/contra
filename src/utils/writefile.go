package utils

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"strconv"
)

// WriteFile saves output to a file
func WriteFile(workspace, config, name string) error {
	// Create file inside workspace folder
	err := workspaceExists(workspace)
	if err == nil {
		err = gitIgnoreExists(workspace)
	}
	if err == nil {
		f, err := os.Create(workspace + "/" + name)
		if err != nil {
			log.Printf("Unable to create config file %s\n", name)
			return err
		}
		defer f.Close()
		// write config data
		f.WriteString(config)
		f.Close()
		f.Sync()
	}
	return err
}

// workspaceExists checks for or creates the workspace directory
func workspaceExists(workspace string) error {
	// Check if the workspace is a directory
	if ws, err := os.Stat(workspace); err == nil && ws.IsDir() {
		return nil
	}

	// Create directory if it isn't there
	err := os.Mkdir(workspace, os.ModePerm)

	return err
}

// gitIgnoreExists checks for or creates .gitignore in the workspace
func gitIgnoreExists(workspace string) error {
	filename := path.Join(workspace, ".gitignore")
	if _, err := os.Stat(filename); err == nil {
		return nil
	}
	ignore := []byte("diffs")
	return ioutil.WriteFile(filename, ignore, 0644)
}

// WriteRunResult will write the count value into the runresult.log file and update the timestamp each run.
func WriteRunResult(filename string, count int) error {
	d1 := []byte(strconv.Itoa(count))
	return ioutil.WriteFile(filename, d1, 0644)
}
