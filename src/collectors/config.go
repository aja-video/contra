package collectors

import (
	"github.com/go-ini/ini"
	"log"
)

// Cyclical import cycle if config is not in collectors...

// ExampleIni is a proof of concept for pulling ini values.
func ExampleIni() string {
	iniFile, err := ini.Load("example.conf")

	if err != nil {
		// Critical...
		log.Println("example.conf")
		panic(err)
	}

	// [neat]  result = awesome
	return iniFile.Section("neat").Key("result").String()
}

// GenerateCollectorsFromConfig is neat
func GenerateCollectorsFromConfig() []Collector {
	// which := []string{"pfsense"}
	//collector := &DevicePfsense{}

	collectorSlice := []Collector{
		&DevicePfsense{},
	}

	// Run sample collectors.
	//CollectpfSense()
	//CollectComware()
	//CollectProcurve()
	//CollectCsb()

	return collectorSlice
}

// FetchConfig is a quick function to pull a config during development.
func FetchConfig(target string) map[string]string {
	iniFile, err := ini.Load("devices.conf")

	if err != nil {
		// Critical...
		log.Println("example.conf")
		panic(err)
	}

	return iniFile.Section(target).KeysHash()
}
