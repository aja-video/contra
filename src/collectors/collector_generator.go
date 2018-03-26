package collectors

// GenerateCollectorsFromConfig is neat
//func GenerateCollectorsFromConfig(config *configuration.Config) []Collector {
//	// which := []string{"pfsense"}
//	//collector := &DevicePfsense{}
//
//	log.Println(config)
//
//	collectorSlice := make([]Collector, 0, len(config.Devices))
//
//	for _, c := range config.Devices {
//		if c.Disabled {
//			// Device config set to disabled=true to quickly turn off.
//			log.Printf("Config disabled: %v", c.Name)
//			continue
//		}
//
//		collector, err := MakeCollector(c)
//
//		if err != nil {
//			log.Printf("Invalid configuration: %v", c.Name)
//			log.Print(err)
//			continue
//		}
//		collectorSlice = append(collectorSlice, collector)
//	}
//
//	// Run sample collectors.
//	//CollectpfSense()
//	//CollectComware()
//	//CollectProcurve()
//	//CollectCsb()
//
//	return collectorSlice
//}

//func GenerateCollectorFromConfig(device configuration.DeviceConfig) Collector {
//	return MakeCollector(device)
//}

// FetchConfig is a quick function to pull a config during development.
//func FetchConfig(target string) map[string]string {
//	iniFile, err := ini.Load("contra.conf")
//
//	if err != nil {
//		// Critical...
//		log.Println("contra.conf")
//		panic(err)
//	}
//
//	return iniFile.Section(target).KeysHash()
//}
