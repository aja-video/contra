package collectors

import (
	"contra/src/utils"
	"fmt"
	"regexp"
	"time"
)

// CollectComware pulls the device config for a comware device.
func CollectComware() string {
	fmt.Printf("Collect Works - Comware\n")

	// set up ssh connection
	s := new(utils.SSHConfig)

	creds := utils.FetchConfig("comware")
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
		1: "<.*.>", // 1 : should always match the initial connection string
		2: "<.*.>",
		3: "return",
	}

	// Commands we will send in response to output above
	send := map[int]string{
		1: "screen-length disable\n",
		2: "display current-configuration\n",
	}

	// Build batcher
	batch := utils.BuildBatcher(send, receive)

	// call GatherExpect to collect the configs
	result, err := utils.GatherExpect(batch, time.Second*10, connection)
	if err != nil {
		panic(err)
	}

	// Strip shell commands, grab only the xml file
	config := regexp.MustCompile(`#[\s\S]*?return`)

	match := config.FindStringSubmatch(result[2].Output)

	utils.WriteFile(match[0], "comware.txt")
	return match[0]
}
