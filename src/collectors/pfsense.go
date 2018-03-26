package collectors

import (
	"contra/src/utils"
	"fmt"
	"github.com/google/goexpect"
	"log"
	"regexp"
	"time"
)

// Repeatedly get the sense that this should just be a function, and that a struct is overkill,
// but the struct provides a very elegant way to wrap and call build batcher and parse result
// in a general collect function.

// DevicePfsense write me.
type DevicePfsense struct {
	*CollectorDefinition
}

// BuildBatcher write me.
func (p *DevicePfsense) BuildBatcher() ([]expect.Batcher, error) {
	return utils.SimpleBatcher([][]string{
		{"option:", "8\n"}, // "option:" should always match the initial connection string
		{".*root", "cat /conf/config.xml\n"},
		{"</pfsense>"},
	})
}

// ParseResult write me.
func (p *DevicePfsense) ParseResult(result string) (string, error) {
	// Strip shell commands, grab only the xml file
	config := regexp.MustCompile(`<\?xml version[\s\S]*?<\/pfsense>`)

	match := config.FindStringSubmatch(result)

	return match[0], nil
}

// ModifySSHConfig just temp.
func (p *DevicePfsense) ModifySSHConfig(config utils.SSHConfig) {
	log.Println("Super speshal.")
}

// CollectpfSense collects a pfSense config.
func CollectpfSense() string {
	fmt.Printf("Collect Works - pfSense\n")

	// set up ssh connection
	s := new(utils.SSHConfig)

	creds := FetchConfig("pfsense")
	// Set up SSHConfig
	s.User = creds["user"]
	s.Password = creds["pass"]
	s.Host = creds["host"] + ":" + creds["port"]

	connection, err := utils.SSHClient(*s)

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
