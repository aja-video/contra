package collectors

import (
	"errors"
	"github.com/aja-video/contra/src/configuration"
	"github.com/aja-video/contra/src/utils"
	contraGit "github.com/aja-video/contra/src/utils/git"
	"github.com/google/goexpect"
	"golang.org/x/crypto/ssh"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
)

// CollectorWorker write me.
type CollectorWorker struct {
	RunConfig *configuration.Config
	factory   collectorFactory
	diffs     []string
}

// RunCollectors runs all collectors
func (cw *CollectorWorker) RunCollectors() {
	cw.registerCollectors()
	// Create a channel with a maximum size of configured concurrency
	queue := make(chan struct{}, cw.RunConfig.Concurrency)
	wg := sync.WaitGroup{}
	for _, device := range cw.RunConfig.Devices {
		// sanity check timeout
		if device.SSHTimeout < time.Second {
			log.Println("WARNING: SSH Timeout should be a minimum of 1s. Disabling device")
			device.Disabled = true
		}
		if device.Disabled {
			log.Printf("Config disabled: %v", device.Name)
			continue
		}
		// Add an element to the queue to limit concurrency
		queue <- struct{}{}
		wg.Add(1)

		// Start collection process for each device
		go func(config configuration.DeviceConfig) {
			defer wg.Done()
			diff, err := cw.Run(config)
			if err != nil {
				log.Printf("Worker resulted in an error: %s\n", err.Error())
			} else if len(diff) > 0 {
				log.Printf("changes found for device %s\n", config.Name)
				cw.diffs = append(cw.diffs, diff)
			}
			// Remove an element from the queue when the collection has finished
			<-queue
		}(device)

	}
	// wait for all collections to finish
	wg.Wait()

	// Currently we are simply notifying that we attempted 4 collections, but we do not
	// reduce this number in the case of an expired ssh timer failure. We may want to.
	// Display the completed collections, and write to run result file.
	log.Printf("Completed collections: %d\n", len(cw.RunConfig.Devices))
	err := utils.WriteRunResult(cw.RunConfig.RunResult, len(cw.RunConfig.Devices))
	if err != nil {
		// If we fail to write the file, we display a note, but do not worry about it here.
		// In theory, a monitoring application (Check MK) will notice the problem.
		log.Printf("Failed to write run result: %s\n", err.Error())
	}
	// send email if anything changed.
	if len(cw.diffs) > 0 {
		// Attempt to send email. Email failure is logged, but does not interrupt this process.
		contraGit.GitSendEmail(cw.RunConfig, cw.diffs)
	}

}

// Run the collector for this device.
func (cw *CollectorWorker) Run(device configuration.DeviceConfig) (string, error) {
	log.Printf("Collect Start: %s\n", device.Name)
	var lastResult string

	// call buildCollector to assemble the collector and batcher
	collector, batchSlice, err := cw.buildCollector(device)
	if err != nil {
		return "", err
	}

	// route53 and http-json do not use ssh
	// TODO: This should be fixed - quick solution to skip SSH for a single collector, but type doesn't scale
	if device.Type != "route53" && device.Type != "http-json" {

		var connection *ssh.Client

		// Set up device SSHConfig
		s := cw.buildSSH(collector, device)

		// handle client timeouts
		clientConnected := make(chan struct{})

		go func(err *error) {
			// set outer scope err
			// this looks like a terrible idea, but the blocking select {} below makes it safe
			connection, *err = utils.SSHClient(*s)
			// close() immediately triggers <-clientConnected
			close(clientConnected)
		}(&err)

		// wait for client
		select {
		case <-clientConnected:
			if err != nil {
				return "", cw.collectFailure(device, err)
			}
			break
		case <-time.After(device.SSHTimeout):
			return "", cw.collectFailure(device, errors.New("SSH Client timeout"))
		}
		// call GatherExpect to collect the configs
		result, err := utils.GatherExpect(batchSlice, device.SSHTimeout, connection)
		if err != nil {
			// Close the connection if collection fails
			closeErr := connection.Close()
			if closeErr != nil {
				log.Printf("WARNING: Unable to close SSH connection for device %s %s\n", device.Name, closeErr.Error())
			}
			// return the collection error
			return "", cw.collectFailure(device, err)
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
		lastResult = result[len(result)-1].Output
	}

	parsed, err := collector.ParseResult(lastResult)
	if err != nil {
		return "", err
	}
	diff, err := contraGit.GitDiff(cw.RunConfig.Workspace, device.Name+".txt", parsed)
	if err != nil {
		log.Printf("error parsing changes in configuration for %s: %s\n", device.Name, err.Error())
	}

	log.Printf("Writing: %s\nLength: %d\n", device.Name, len(parsed))

	err = utils.WriteFile(cw.RunConfig.Workspace, parsed, device.Name+".txt")

	return diff, err

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
		_ = utils.SendEmail(cw.RunConfig, "Contra failure warning!", strings.Join(message, " "))
		// Empty the channel after we've reset a notification
		for len(d.FailChan) > 0 {
			<-d.FailChan
		}
	}

	// log a warning, reminder the matcher could be out of date, or incomplete.
	log.Printf("WARNING: Contra failed to gather configs from %s with error: %s\n"+
		"This can happen if the incorrect device type is selected, or if the device sends "+
		"output that does not match the expected output.", d.Name, err.Error())

	// We're logging it already, don't pass it further.
	//return err
	return nil
}

func (cw *CollectorWorker) registerCollectors() {
	for name, device := range deviceMap {
		cw.factory.Register(name, device)
	}
}

// buildCollector puts the pieces together for a working collector
func (cw *CollectorWorker) buildCollector(device configuration.DeviceConfig) (Collector, []expect.Batcher, error) {
	// Initialize Collector
	collector, err := cw.factory.MakeCollector(device.Type)
	if err != nil {
		return nil, nil, err
	}
	collector.SetDeviceConfig(device)

	batchSlice, err := collector.BuildBatcher()
	if err != nil {
		return nil, nil, err
	}
	return collector, batchSlice, err
}

func (cw *CollectorWorker) buildSSH(collector Collector, device configuration.DeviceConfig) *utils.SSHConfig {
	s := &utils.SSHConfig{
		User:          device.User,
		Pass:          device.Pass,
		Host:          device.Host + ":" + strconv.Itoa(device.Port),
		AuthMethod:    device.SSHAuthMethod,
		PrivateKey:    device.SSHPrivateKey,
		AllowInsecure: device.AllowInsecureSSH,
		SSHTimeout:    device.SSHTimeout,
	}

	// Special case... only some collectors need to make some modifications.
	if collectorSpecial, ok := collector.(CollectorSpecialSSH); ok {
		collectorSpecial.ModifySSHConfig(s)
	}
	if collectorSpecial, ok := collector.(CollectorSpecialTerminal); ok {
		collectorSpecial.ModifyUsername(s)
	}

	// Pull in device config Cipher overrides if necessary.
	if device.Ciphers != "" {
		s.Ciphers = strings.Split(device.Ciphers, ",")
	}
	return s
}
