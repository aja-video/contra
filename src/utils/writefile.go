package utils

import (
	"contra/src/configuration"
	"log"
	"os"
)

// WriteFile saves output to a file
func WriteFile(c configuration.Config, config string, name string) error {
	// Create file inside workspace folder
	err := wsDir(c)
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

// wsDir checks for or creates the workspace directory
func wsDir(c configuration.Config) error {
	// Check if the workspace is a directory
	if ws, err := os.Stat(c.Workspace); err == nil && ws.IsDir() {
		return nil
	}

	// Create directory if it isn't there
	err := os.Mkdir(c.Workspace, os.ModePerm)

	return err
}
