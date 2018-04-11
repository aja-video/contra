package collectors

import (
	"contra/src/configuration"
	"contra/src/utils"
	"log"
	"strconv"
	"strings"
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

	log.Printf("Completed collections: %d\n", len(cw.RunConfig.Devices))
}

// Run the collector for this device.
func (cw *CollectorWorker) Run(device configuration.DeviceConfig) error {
	log.Printf("Collect Start: %s\n", device.Name)

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

	// Pull in device config Cipher overrides if necessary.
	if device.Ciphers != "" {
		s.Ciphers = strings.Split(device.Ciphers, ",")
	}

	connection, err := utils.SSHClient(*s)
	if err != nil {
		return cw.collectFailure(device, err)
	}
	// call GatherExpect to collect the configs
	result, err := utils.GatherExpect(batchSlice, time.Second*10, connection)
	if err != nil {
		// Close the connection if collection fails
		closeErr := connection.Close()
		if closeErr != nil {
			log.Printf("WARNING: Unable to close SSH connection for device %s %s\n", device.Name, closeErr.Error())
		}
		// return the collection error
		return cw.collectFailure(device, err)
	}
	// Read from FailChan if it isn't empty
	if len(device.FailChan) > 0 {
		<-device.FailChan
	}
	// Close ssh connection
	err = connection.Close()
	if err != nil {
		log.Printf("WARNING: Error closing SSH Connection for %s: %s\n", device.Name, err.Error())
	}
	// Grab just the last result.
	lastResult := result[len(result)-1].Output
	parsed, _ := collector.ParseResult(lastResult)

	log.Printf("Writing: %s\nLength: %d\n", device.Name, len(parsed))

	utils.WriteFile(*cw.RunConfig, parsed, device.Name+".txt")

	return nil
}

// collectFailure handles collector failures
func (cw *CollectorWorker) collectFailure(d configuration.DeviceConfig, err error) error {
	// define email notification content
	var message []string
	message = append(message, "Contra was unable to gather configs from",
		d.Name, strconv.Itoa(d.FailureWarning), "times", "last error:", err.Error())
	// Add an element to the warning queue if it isn't full
	if len(d.FailChan) < cap(d.FailChan) {
		d.FailChan <- true
	} else if cap(d.FailChan) > 0 {
		// Fire off an email if the warning queue is full and email is enabled
		utils.SendEmail(cw.RunConfig, "Contra failure warning!", strings.Join(message, " "))
		// Empty the channel after we've reset a notification
		for len(d.FailChan) > 0 {
			<-d.FailChan
		}
	}
	// log a warning
	log.Printf("WARNING: Contra failed to gather configs from %s with error: %s\n", d.Name, err.Error())

	return err
}
