package main

import (
	"contra/src/collectors"
	"contra/src/utils"
	"fmt"
)

func main() {
	// Print something.
	fmt.Printf("Contra\n")

	// Print something from an imported package.
	collectors.Collect()
	collectors.CollectComware()

	// Pull a result from an ini using a third party package
	str := utils.ExampleIni()
	fmt.Printf("Found: %v\n", str)

}
