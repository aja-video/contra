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

	// set up ssh connection
	s := new(utils.SSHConfig)

	creds := utils.FetchConfig("pfsense")
	// Set up SSHConfig
	s.User = creds["user"]
	s.Password = creds["pass"]
	s.Host = creds["host"] + ":" + creds["port"]

	connection, err := utils.SSHClient(*s)

	if err != nil {
		panic(err)
	}

	// Output we expect to receive
	receive := map[int]string{
		1: "option:", // 1 : should always match the initial connection string
		2: ".*root",
		3: "</pfsense>",
	}

	// Commands we will send in response to output above
	send := map[int]string{
		1: "8\n",
		2: "cat /conf/config.xml\n",
	}

	// Build batcher
	batch := utils.BuildBatcher(send, receive)

	// call GatherExpect to collect the configs
	result, err := utils.GatherExpect(batch, time.Second*10, connection)
	if err != nil {
		panic(err)
	}

	// Strip shell commands, grab only the xml file
	config := regexp.MustCompile(`<\?xml version[\s\S]*?<\/pfsense>`)

	match := config.FindStringSubmatch(result[2].Output)

	utils.WriteFile(match[0], "pfsense.txt")
	return match[0]
}
