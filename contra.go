package main

import (
	"contra/src/collectors"
	"contra/src/utils"
	"fmt"
)

func main() {
	// Print something.
	fmt.Printf("Contra\n")

	// Run sample collectors.
	//collectors.Collect()
	//collectors.CollectComware()
	//collectors.CollectProcurve()
	collectors.CollectCsb()
	// Pull a result from an ini using a third party package
	str := utils.ExampleIni()
	fmt.Printf("Found: %v\n", str)

	utils.Commit()
}
