package core

import (
	"contra/src/collectors"
	"contra/src/configuration"
	"contra/src/utils/git"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"
)

const version = "1.0.0"

// Application holds global application data and functions for kicking off execution.
type Application struct {
	config *configuration.Config
}

// Start is the main entrance to the application.
func (a *Application) Start() {
	// Parse the config, which brings in flags.
	a.config = configuration.GetConfig()

	// Display banner second, in case config declares a quiet run.
	a.DisplayBanner()

	// TODO: This would probably be better to only run ONCE on StandardRun/Daemonize.
	// Now that we have completely determined our configs (including command line flags)
	// If we want to encrypt passwords, then kick it off.
	if a.config.EncryptPasswords {
		configuration.EncryptConfigFile(a.config.ConfigFile)
	}

	// Determine what to do.
	a.Route()
}

//RunDaemon will persist and run collectors at the configured interval
func (a *Application) RunDaemon() {
	interval := a.config.Interval
	for {
		a.StandardRun()
		log.Printf("Collection finished, sleeping for %s\n", interval)
		time.Sleep(interval)
	}
}

// StandardRun if there are no special cases designated by the configuration.
func (a *Application) StandardRun() error {
	// Initialize our main worker.
	worker := collectors.CollectorWorker{
		RunConfig: a.config,
	}

	// Collect everything
	worker.RunCollectors()

	// And check for any necessary commits.
	err := utils.GitOps(a.config)
	if err != nil {
		return err
	}

	return nil
}

// Route determines what to do, and kicks off the doing.
func (a *Application) Route() {
	// Determine if the config designates some special run process, otherwise handle our main handler.
	if a.config.Copyrights {
		a.DisplayCopyrights()
	} else if a.config.Debug {
		a.DisplayDebugInfo()
	} else if a.config.Version {
		a.DisplayVersion()
	} else if a.config.Daemonize {
		// Repeat collectors every interval
		a.RunDaemon()
	} else {
		// Standard operating procedure.
		log.Println("Contra is not configured to run as a Daemon, performing a single collection")
		a.StandardRun()
	}
}

// DisplayBanner with basic information about this application.
func (a *Application) DisplayBanner() {
	// Print something.
	if !a.config.Quiet {
		fmt.Printf("\n=== Contra ===\n"+
			" - Network Device Configuration Tracking\n"+
			" - AJA Video Systems, Inc. Version: %s\n\n", version)
		log.Printf("Contra started with configuration file %s\n", a.config.ConfigFile)
	}
}

// DisplayCopyrights simply dumping the COPYRIGHTS file.
func (a *Application) DisplayCopyrights() {
	log.Println("COPYRIGHT Information")
	data, err := ioutil.ReadFile("COPYRIGHTS")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	log.Println(string(data))
	os.Exit(0)
}

// DisplayDebugInfo may print sensitive passwords to the screen.
func (a *Application) DisplayDebugInfo() {
	log.Println("DEBUG ENABLED: Dumping config and exiting.")
	log.Println(a.config)
	os.Exit(0)
}

// DisplayVersion prints the Contra version and exits
func (a *Application) DisplayVersion() {
	// If Quiet is set, just display the version.
	// If Quiet is not set, the version is included in from DisplayBanner.
	if a.config.Quiet {
		fmt.Println(version)
	}

	os.Exit(0)
}
