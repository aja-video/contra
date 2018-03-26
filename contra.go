package main

import (
	"contra/src/collectors"
	"contra/src/configuration"
	"contra/src/utils"
	"fmt"
)

func main() {
	// Print something.
	fmt.Printf("\n=== Contra ===\n - Network Device Configuration Tracking\n - AJA Video Systems\n\n")

	config := configuration.GetConfig()

	worker := collectors.CollectorWorker{
		RunConfig: config,
	}

	worker.RunCollectors()

	utils.Commit()
}
