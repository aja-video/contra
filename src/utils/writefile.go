package utils

import (
	"contra/src/configuration"
	"fmt"
	"os"
)

// WriteFile saves output to a file
func WriteFile(c configuration.Config, config string, name string) error {
	fmt.Println(c.Workspace)
	// Create file inside workspace folder
	f, err := os.Create(c.Workspace + name)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	// write config data
	f.WriteString(config)
	f.Close()

	return err
}
