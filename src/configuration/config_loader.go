package configuration

import (
	"log"
	"sync"
	"time"
)

// ConfigLoader - Orchestrates the configuration loading.
// Config - The plain config type.
// ConfigFile - Loads configuration data from a file.
// ConfigFlags - Loads configuration data from a flag.

// Singleton pattern ensures a single config across concurrent threads.
var instance *Config
var once sync.Once

// GetConfig is concurrency safe loading and retrieving of the config data.
func GetConfig() *Config {
	once.Do(func() {
		instance = loadConfig()
	})
	return instance
}

// loadConfig fetches the various sources of configuration data, and returns the fully prepared config.
func loadConfig() *Config {
	// Parse command line flags
	parseConfigFlags()

	// Pull a copy of the defaults, and convert to a pointer.
	config := getConfigDefaults()

	// Load config data from INI file on top of default values.
	configPath := configFlagsGetConfigPath()
	mergeConfigFile(config, configPath)

	// Fetch flags and merge on top of file+default values.
	mergeConfigFlags(config)

	// Sanity Check
	if config.Interval < time.Second {
		// TODO: Friendly way to exit?
		log.Fatalln("Interval should be a minimum of 1s. Did you forget the seconds?")
	}
	if config.Timeout < time.Second {
		// TODO: Friendly way to exit?
		log.Fatalln("Timeout should be a minimum of 1s. Did you forget the seconds?")
	}

	return config
}
