package main

import (
	"contra/src/collectors"
	"contra/src/utils"
	"fmt"
)

func main() {
	// Print something.
	fmt.Printf("Contra\n")

	worker := collectors.CollectorWorker{}
	collectorSlice := collectors.GenerateCollectorsFromConfig()
	for _, collector := range collectorSlice {
		worker.Run(collector)
	}

	// Pull a result from an ini using a third party package
	str := collectors.ExampleIni()
	fmt.Printf("Found: %v\n", str)

	utils.Commit()
}
