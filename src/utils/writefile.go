package utils

import "os"

// WriteFile saves output to a file
func WriteFile(config string, name string) error {

	// Create file inside workspace folder
	f, err := os.Create("workspace/" + name)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	// write config data
	f.WriteString(config)
	f.Close()

	return err
}
