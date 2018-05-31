package utils

import (
	"contra/src/configuration"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

// WriteFile saves output to a file
func WriteFile(c configuration.Config, config string, name string) error {
	// Create file inside workspace folder
	err := workspaceExists(c)
	if err == nil {
		f, err := os.Create(c.Workspace + "/" + name)
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
func workspaceExists(c configuration.Config) error {
	// Check if the workspace is a directory
	if ws, err := os.Stat(c.Workspace); err == nil && ws.IsDir() {
		return nil
	}

	// Create directory if it isn't there
	err := os.Mkdir(c.Workspace, os.ModePerm)

	return err
}

// WriteRunResult will write the count value into the runresult.log file and update the timestamp each run.
func WriteRunResult(count int) error {
	name := "runresult.log"
	d1 := []byte(strconv.Itoa(count))
	return ioutil.WriteFile(name, d1, 0644)
}
