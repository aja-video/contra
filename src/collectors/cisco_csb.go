package collectors

import (
	"contra/src/utils"
	"fmt"
	"regexp"
	"time"
)

// CollectCsb pulls the device config for a Cisco Small Business device.
func CollectCsb() string {
	fmt.Printf("Collect Works - Cisco_csb\n")

	// set up ssh connection
	s := new(utils.SSHConfig)

	creds := utils.FetchConfig("csb")
	// Set up SSHConfig
	s.User = creds["user"]
	s.Password = creds["pass"]
	s.Host = creds["host"] + ":" + creds["port"]
	s.Ciphers = []string{"aes256-cbc", "aes128-cbc"}

	connection, err := utils.SSHClient(*s)

	if err != nil {
		panic(err)
	}

	// Output we expect to receive
	receive := map[int]string{
		1: "Name:", // This is assuming prompt for User Name on Cisco CSB - this may not always be the case
		2: ".*#",
		3: ".*#",
		4: ".*#",
	}

	// Commands we will send in response to output above
	send := map[int]string{
		1: creds["user"] + "\n",
		2: "terminal datadump\n",
		3: "show running-config",
	}

	// Build batcher
	batch := utils.BuildBatcher(send, receive)

	// call GatherExpect to collect the configs
	result, err := utils.GatherExpect(batch, time.Second*10, connection)
	if err != nil {
		panic(err)
	}

	// Strip shell commands, grab only the xml file
	config := regexp.MustCompile(`config-file-header[\s\S]*?#`) // This may break if there is a '#' in the config

	match := config.FindStringSubmatch(result[3].Output)

	//utils.WriteFile(match[0], "cisco_csb.txt")
	fmt.Printf(match[0])
	return match[0]
}
