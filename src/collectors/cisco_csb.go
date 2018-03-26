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
		1: "User Name:", // This is assuming prompt for User Name on Cisco CSB - this may not always be the case
		2: "Password:",
		3: ".*#",
		4: ".*#",
		5: ".*#",
	}

	// Commands we will send in response to output above
	send := map[int]string{
		1: creds["user"] + "\n",
		2: creds["pass"] + "\n",
		3: "terminal datadump\n",
		4: "show running-config\n",
	}

	// Build batcher
	batch := utils.BuildBatcher(send, receive)

	// call GatherExpect to collect the configs
	result, err := utils.GatherExpect(batch, time.Second*10, connection)
	if err != nil {
		panic(err)
	}

	fmt.Print(len(result))
	// Strip shell commands, grab only the xml file
	config := regexp.MustCompile(`config-file-header[\s\S]*?#`) // This may break if there is a '#' in the config

	match := config.FindStringSubmatch(result[len(result)-1].Output)

	utils.WriteFile(match[0], "cisco_csb.txt")
	fmt.Printf(match[0])
	return match[0]
}
