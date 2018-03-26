package collectors

import (
	"contra/src/utils"
	"fmt"
	"regexp"
	"time"
)

// Collect currently collects a pfSense config.
func Collect() string {
	fmt.Printf("Collect Works - pfSense\n")

	// set up ssh connection - obviously not the right place for this
	creds := utils.FetchConfig("pfsense")
	connection, err := utils.SSHClient(creds["user"], creds["pass"], creds["host"]+":"+creds["port"])

	if err != nil {
		panic(err)
	}

	batch, err := utils.SimpleBatcher([][]string{
		{"option:", "8\n"}, // "option:" should always match the initial connection string
		{".*root", "cat /conf/config.xml\n"},
		{"</pfsense>"},
	})

	if err != nil {
		panic(err)
	}

	// call GatherExpect to collect the configs
	// TODO: Verify pointer/reference/dereference is necessary.
	result, err := utils.GatherExpect(&batch, time.Second*10, connection)
	if err != nil {
		panic(err)
	}

	// Strip shell commands, grab only the xml file
	config := regexp.MustCompile(`<\?xml version[\s\S]*?<\/pfsense>`)

	match := config.FindStringSubmatch(result[2].Output)

	utils.WriteFile(match[0], "pfsense.txt")
	return match[0]
}
