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
	creds := utils.FetchConfig("pfsense")
	connection, _ := utils.SSHClient(creds["user"], creds["pass"], creds["host"]+":"+creds["port"])
	result, err := utils.GatherExpect(&batchSlice, time.Second*10, connection)
	if err != nil {
		panic(err)
	}

	// TODO: Grab just the last result.

	parsed, _ := collector.ParseResult(result[2].Output)

	// result[2].Output
	// match[0]
	log.Printf("Writing: %s\n%s\n", which, parsed)

	utils.WriteFile(parsed, which+".txt")
}
