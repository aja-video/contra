package collectors

import (
	"contra/src/configuration"
	"contra/src/utils"
	"fmt"
	"log"
	"strconv"
	"time"
)

// CollectorWorker write me.
type CollectorWorker struct {
	RunConfig *configuration.Config
}

// RunCollectors runs all collectors
func (cw *CollectorWorker) RunCollectors() {
	// Create a channel with a maximum size of configured concurrency
	queue := make(chan bool, cw.RunConfig.Concurrency)
	for _, device := range cw.RunConfig.Devices {
		if device.Disabled {
			log.Printf("Config disabled: %v", device.Name)
			continue
		}
		// Add an element to the queue for each enabled device
		queue <- true
		// Start collection process for each device
		go func(config configuration.DeviceConfig) {
			cw.Run(config)
			// Remove an element from the queue when the collection has finished
			defer func() {
				<-queue
			}()
		}(device)

	}
	// If we can add to the queue a number of elements equal to concurrency
	// our goroutines are finished and we can leave this function
	for l := 0; l < cap(queue); l++ {
		queue <- true
	}

	fmt.Printf("Completed collections: %d\n", len(cw.RunConfig.Devices))
}

// Run the collector for this device.
func (cw *CollectorWorker) Run(device configuration.DeviceConfig) {
	fmt.Printf("Collect Start: %s\n", device.Name)

	collector, _ := MakeCollector(device)

	batchSlice, _ := collector.BuildBatcher()

	// Set up SSHConfig
	s := &utils.SSHConfig{
		User: device.User,
		Pass: device.Pass,
		Host: device.Host + ":" + strconv.Itoa(device.Port),
	}

	// Special case... only some collectors need to make some modifications.
	if collectorSpecial, ok := collector.(CollectorSpecial); ok {
		collectorSpecial.ModifySSHConfig(s)
	}

	connection, err := utils.SSHClient(*s)

	// call GatherExpect to collect the configs
	// TODO: Verify pointer/reference/dereference is necessary.
	result, err := utils.GatherExpect(batchSlice, time.Second*10, connection)
	if err != nil {
		panic(err)
	}

	// Grab just the last result.
	lastResult := (result)[len((result))-1].Output
	parsed, _ := collector.ParseResult(lastResult)

	log.Printf("Writing: %s\nLength: %d\n", device.Name, len(parsed))

	utils.WriteFile(*cw.RunConfig, parsed, device.Name+".txt")

}
