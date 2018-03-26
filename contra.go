package main

import (
	"contra/src/collectors"
	"contra/src/utils"
	"fmt"
)

func main() {
	// Print something.
	fmt.Printf("Contra\n")

	// TODO: rework this... config parsing should build a struct of collectors.
	which := []string{"pfsense"}
	for range which {
		collector := &collectors.DevicePfsense{}
		worker := collectors.CollectorWorker{}
		worker.Run(collector)
	}

	// Run sample collectors.
	collectors.CollectpfSense()
	collectors.CollectComware()
	collectors.CollectProcurve()
	collectors.CollectCsb()

	// Pull a result from an ini using a third party package
	str := utils.ExampleIni()
	fmt.Printf("Found: %v\n", str)

	utils.Commit()
}
