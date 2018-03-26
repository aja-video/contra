package collectors

import (
	"contra/src/configuration"
	"contra/src/utils"
	"fmt"
	"log"
	"strconv"
	"time"
)

// TODO: Probably rename this to collector, and rename collector to CollectorDefinition

// CollectorWorker write me.
type CollectorWorker struct {
	RunConfig *configuration.Config
}

// RunCollectors runs all collectors
func (cw *CollectorWorker) RunCollectors() {
	//collectorSlice := GenerateCollectorsFromConfig(cw.RunConfig)
	for _, device := range cw.RunConfig.Devices {
		if device.Disabled {
			log.Printf("Config disabled: %v", device.Name)
			continue
		}

		cw.Run(device)
	}

	fmt.Printf("Completed collections: %d\n", len(cw.RunConfig.Devices))
}

// Run write me.
func (cw *CollectorWorker) Run(device configuration.DeviceConfig) {
	fmt.Printf("Collect Start: %s\n", device.Name)

	collector, _ := MakeCollector(device)

	batchSlice, _ := collector.BuildBatcher()

	log.Print(batchSlice)

	// call GatherExpect to collect the configs
	// set up ssh connection - obviously not the right place for this
	//creds := FetchConfig("pfsense")

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

	result, err := utils.GatherExpect(&batchSlice, time.Second*10, connection)
	if err != nil {
		panic(err)
	}

	// DONE: Grab just the last result.
	// result[2].Output
	lastResult := result[len(result)-1].Output
	parsed, _ := collector.ParseResult(lastResult)

	// match[0]
	log.Printf("Writing: %s\n%d\n", device.Name, len(parsed))

	utils.WriteFile(parsed, device.Name+".txt")
}
