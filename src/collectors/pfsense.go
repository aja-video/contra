package collectors

import (
	"contra/src/utils"
	"fmt"
	"github.com/google/goexpect"
	"regexp"
	"time"
)

func Collect() string {
	fmt.Printf("Collect Works\n")

	// set up ssh connection - obviously not the right place for this
	connection, err := utils.SshClient("admin", "thisshouldn'tbehere!", "192.168.1.1:22")

	if err != nil {
		panic(err)
	}

	// Set up expect.Batcher with commands and expected output
	pf := []expect.Batcher{
		&expect.BExp{
			R: `option:`,
		},
		&expect.BSnd{
			S: "8\n",
		},
		&expect.BExp{
			R: ".*root",
		},
		&expect.BSnd{
			S: "cat /conf/config.xml\n",
		},
		&expect.BExp{
			R: "</pfsense>",
		},
	}

	// call GatherExpect to collect the configs
	result, err := utils.GatherExpect(&pf, time.Second*10, connection)
	if err != nil {
		panic(err)
	}

	// Strip shell commands, grab only the xml file
	config := regexp.MustCompile(`<\?xml version[\s\S]*?<\/pfsense>`)

	match := config.FindStringSubmatch(result[2].Output)

	utils.WriteFile(match[0], "pfsense.txt")
	return match[0]
}
