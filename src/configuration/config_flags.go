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
	//flag.Bool("v", true, "Display additional output details.")

	// Add reminder that we're not currently using this.
	flag.Lookup("v").Usage = "log level for V logs (not currently used)"

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
// some other edge case considerations. The bool config map is to
// reduce the cyclic complexity slightly. Wish the flags library would
// have something built in, perhaps need to write a mapping class.
func mergeConfigFlags(config *Config) {

	boolConfigMap := map[string]*bool{
		"d":          &config.Daemonize,
		"dc":         &config.DisableCollection,
		"debug":      &config.Debug,
		"version":    &config.Version,
		"x":          &config.AllowInsecureSSH,
		"q":          &config.Quiet,
		"e":          &config.EmailEnabled,
		"w":          &config.WebserverEnabled,
		"copyrights": &config.Copyrights,
	}

	flag.Visit(func(flagVal *flag.Flag) {
		if val, ok := boolConfigMap[flagVal.Name]; ok {
			// We found the value, dereference and assign.
			*val = flagVal.Value.(flag.Getter).Get().(bool)
			return
		}

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
		case "listen":
			config.HTTPListen = flagVal.Value.(flag.Getter).Get().(string)
		default:
			// Fail if not defined.
			log.Fatalf("Flag merge not configured for %v", flagVal)
		}
	})
}

func configFlagsGetConfigPath() string {
	// Flags must be parsed.
	if !flag.Parsed() {
		panic("Flags not yet parsed.")
	}

	// Grab the default config path.
	configPath := getConfigDefaults().ConfigFile

	// And check to see if we want to override it with a flag.
	configFlag := flag.Lookup("c")
	if configFlag != nil {
		found := flag.Lookup("c").Value.(flag.Getter).Get().(string)
		if found != configPath {
			//log.Printf("Config: %s, Switching to found: %s", configPath, found)
			return found
		}
	}
	// if the default config exists use it
	if _, err := os.Stat(configPath); err == nil {
		return configPath
	}
	// try for /etc/contra.conf
	if _, err := os.Stat(`/etc/contra.conf`); err == nil {
		return `/etc/contra.conf`
	}
	// Die with a useful error if we can't find a config file
	log.Fatalf("ERROR: Unable to open config file %s", configPath)

	// Nothing to return here
	return ""
}
