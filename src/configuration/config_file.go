package configuration

import (
	"github.com/go-ini/ini"
	"log"
)

func mergeConfigFile(config *Config, filePath string) {
	iniFile, err := ini.Load(filePath)
	if err != nil {
		log.Println(filePath)
		panic(err) // Critical....
		//return nil, err
	}

	// Map [main] to config.
	err = iniFile.Section("main").MapTo(config)
	if err != nil {
		panic(err)
	}

	// Map/Load Device Configs
	for _, section := range iniFile.Sections() {
		if section.Name() == "main" || section.Name() == "DEFAULT" {
			continue
		}
		if !section.HasKey("Type") {
			log.Panicf("Device [%v] must have a type defined.", section.Name())
		}
		deviceConfig := DeviceConfig{
			FailureWarning: 5,
		}
		section.MapTo(&deviceConfig)
		// Copy the section name into the device config for reference.
		deviceConfig.Name = section.Name()
		// Set up device failure channel if it isn't disabled (zero value)
		if deviceConfig.FailureWarning > 0 {
			deviceConfig.FailChan = make(chan bool, deviceConfig.FailureWarning)
		}
		config.Devices = append(config.Devices, deviceConfig)
	}
}
