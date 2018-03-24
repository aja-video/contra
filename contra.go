package main

import (
	"contra/src/collectors"
	"contra/src/utils"
	"fmt"
	"time"
)

func main() {
	// Print something.
	fmt.Printf("Contra\n")

	// Print something from an imported package.
	collectors.Collect()

	// Pull a result from an ini using a third party package
	str := utils.ExampleIni()
	fmt.Printf("Found: %v\n", str)

	// Call expect with command, timeout, and desired result
	_, match, e := utils.Expect("ls", 10*time.Second, "contra.go")
	if e != nil {
		panic(e)
	}

	// Print expect match
	fmt.Print("Expect regex found file: ", match[0], "\n")

}
