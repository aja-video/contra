package configuration

import (
	"flag"
	"log"
	"os"
	"time"
)

func parseConfigFlags() {
	defaults := getConfigDefaults()

	// TODO: "v" is declared by glog, do we want to use that?
	//flag.Bool("v", defaults.Verbose, "Display additional output details.")

	flag.String("c", defaults.ConfigFile, "Config file name.")
	flag.Int("p", defaults.Concurrency, "Parallel concurrent threads to use for collection.")
	flag.Duration("i", defaults.Interval, "Interval in seconds between run calls.")
	flag.Duration("t", defaults.Timeout, "Timeout default in seconds to wait for collection to finish.")
	flag.Bool("q", defaults.Quiet, "Suppress most output except for problems or warnings.")
	flag.Bool("debug", defaults.Debug, "Enable DEBUG flag for development.")
	flag.Bool("x", defaults.AllowInsecureSSH, "Allow untrusted SSH keys.")
	flag.Bool("e", defaults.EmailEnabled, "Enable or disable email when changes found.")
	flag.Bool("dc", defaults.DisableCollection, "Disable collector processing.")
	flag.Bool("w", defaults.WebserverEnabled, "Run a web status server.")
	flag.String("listen", defaults.HTTPListen, "Host and port to use for HTTP status server.")
	flag.Bool("copyrights", defaults.Copyrights, "Display copyright licenses of compiled packages.")
	flag.Bool("d", defaults.Daemonize, "Run as Daemon")
	flag.Bool("version", defaults.Version, "Display Contra version")
	flag.Parse()
}

// mergeConfigFlags maps the flag values back onto Config.
// There ought to be a more efficient way to handle this when
// combined with the above function for defining and parsing
// the flags. However, it is not apparent how to easily tell whether
// a flag was explicitly set to a default value or not. Plus,
// some other edge case considerations.
func mergeConfigFlags(config *Config) {
	flag.Visit(func(flagVal *flag.Flag) {
		switch flagVal.Name {
		// These should be in the same order that the flags above are declared.
		case "c":
			config.ConfigFile = flagVal.Value.(flag.Getter).Get().(string)
		case "p":
			config.Concurrency = flagVal.Value.(flag.Getter).Get().(int)
		case "i":
			config.Interval = flagVal.Value.(flag.Getter).Get().(time.Duration)
		case "t":
			config.Timeout = flagVal.Value.(flag.Getter).Get().(time.Duration)
		case "debug":
			config.Debug = flagVal.Value.(flag.Getter).Get().(bool)
		case "dc":
			config.DisableCollection = flagVal.Value.(flag.Getter).Get().(bool)
		case "x":
			config.AllowInsecureSSH = flagVal.Value.(flag.Getter).Get().(bool)
		case "q":
			config.Quiet = flagVal.Value.(flag.Getter).Get().(bool)
		case "e":
			config.EmailEnabled = flagVal.Value.(flag.Getter).Get().(bool)
		case "w":
			config.WebserverEnabled = flagVal.Value.(flag.Getter).Get().(bool)
		case "listen":
			config.HTTPListen = flagVal.Value.(flag.Getter).Get().(string)
		case "copyrights":
			config.Copyrights = flagVal.Value.(flag.Getter).Get().(bool)
		case "d":
			config.Daemonize = flagVal.Value.(flag.Getter).Get().(bool)
			// Fail if not defined.
		case "version":
			config.Version = flagVal.Value.(flag.Getter).Get().(bool)
		default:
			log.Fatalf("Flag merge not configured for %v", flagVal)
		}
	})
}

func configFlagsGetConfigPath() string {
	var configPath string
	// Flags must be parsed.
	if !flag.Parsed() {
		panic("Flags not yet parsed.")
	}

	// Grab the default config path.
	configDefaultPath := getConfigDefaults().ConfigFile

	// And check to see if we want to override it with a flag.
	configFlag := flag.Lookup("c")
	if configFlag != nil {
		found := flag.Lookup("c").Value.(flag.Getter).Get().(string)
		if found != configDefaultPath {
			//log.Printf("Config: %s, Switching to found: %s", configPath, found)
			return found
		}
	}
	// if the default config exists use it
	if _, err := os.Stat(configDefaultPath); err == nil {
		configPath = configDefaultPath
	}
	// try for /etc/contra.conf
	if _, err := os.Stat(`/etc/contra.conf`); err == nil {
		configPath = `/etc/contra.conf`
	}
	// Die with a useful error if we can't find a config file
	if len(configPath) == 0 {
		log.Fatalf("ERROR: Unable to open config file %s", configPath)
	}

	return configPath
}
