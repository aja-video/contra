package configuration

import (
	"log"
	"os"
	"sync"
	"time"
)

// ConfigLoader - Orchestrates the configuration loading.
// Config - The plain config type.
// ConfigFile - Loads configuration data from a file.
// ConfigFlags - Loads configuration data from a flag.

// Singleton pattern ensures a single config across concurrent threads.
var instance *Config
var loadOnce sync.Once
var parseOnce sync.Once

// GetConfig is concurrency safe loading and retrieving of the config data.
func GetConfig() *Config {
	loadOnce.Do(func() {
		instance = loadConfig()
	})
	return instance
}

// ReloadConfig pulls the config in again, useful for making changes to a running service.
func ReloadConfig() {
	instance = loadConfig()
}

// loadConfig fetches the various sources of configuration data, and returns the fully prepared config.
func loadConfig() *Config {
	// Parse command line flags
	parseOnce.Do(func() {
		// We only want to parse (and define) the flags once.
		parseConfigFlags()
	})

	// Pull a copy of the defaults, and convert to a pointer.
	config := getConfigDefaults()

	// Load config data from INI file on top of default values.
	configPath := configFlagsGetConfigPath()

	// Remember the file we picked, such as /etc/contra.conf
	config.ConfigFile = configPath

	// Merge the config file params into the config.
	if err := mergeConfigFile(config, configPath); err != nil {
		log.Fatalf("Error loading config file: %s, %s", configPath, err.Error())
	}

	// Fetch flags and merge on top of file+default values.
	mergeConfigFlags(config)

	// Check for config test
	if err := configTester(config.ConfigFile); err != nil {
		log.Fatalf("Contra configuration test error: %s", err.Error())
	} else {
		if !config.Quiet {
			log.Println("Contra configuration test passed")
		}
		if config.ConfigTest {
			os.Exit(0)
		}
	}

	// Sanity Check
	if config.Interval < time.Second {
		// TODO: Friendly way to exit?
		log.Fatalln("Interval should be a minimum of 1s. Did you forget the seconds?")
	}
	if config.Timeout < time.Second {
		// TODO: Friendly way to exit?
		log.Fatalln("Timeout should be a minimum of 1s. Did you forget the seconds?")
	}

	// Check and decrypt any passwords to the loaded in-memory config.
	if err := decryptLoadedConfig(config); err != nil {
		log.Fatalf("error decrypting configuration: %s", err.Error())
	}

	return config
}
