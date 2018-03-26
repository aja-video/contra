package collectors

import (
	"contra/src/utils"
	"fmt"
	"log"
	"time"
)

// TODO: Probably rename this to collector, and rename collector to CollectorDefinition

// CollectorWorker write me.
type CollectorWorker struct {
}

// Run write me.
func (cw *CollectorWorker) Run(collector Collector) {
	which := "pfsense"

	fmt.Printf("Collect Works - pfSense\n")

	batchSlice, _ := collector.BuildBatcher()

	log.Print(batchSlice)

	// call GatherExpect to collect the configs
	// set up ssh connection - obviously not the right place for this
	creds := FetchConfig("pfsense")

	// set up ssh connection
	s := new(utils.SSHConfig)
	// Set up SSHConfig
	s.User = creds["user"]
	s.Password = creds["pass"]
	s.Host = creds["host"] + ":" + creds["port"]

	// Special case... only some collectors need to make some modifications.
	if collectorSpecial, ok := collector.(CollectorSpecial); ok {
		collectorSpecial.ModifySSHConfig(*s)
	}

	connection, err := utils.SSHClient(*s)

	result, err := utils.GatherExpect(&batchSlice, time.Second*10, connection)
	if err != nil {
		panic(err)
	}

	// TODO: Grab just the last result.

	parsed, _ := collector.ParseResult(result[2].Output)

	// result[2].Output
	// match[0]
	log.Printf("Writing: %s\n%d\n", which, len(parsed))

	utils.WriteFile(parsed, which+".txt")
}
