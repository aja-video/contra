package configuration

import (
	"fmt"
	"github.com/go-ini/ini"
)

func mergeConfigFile(config *Config, filePath string) error {
	iniFile, err := ini.Load(filePath)
	if err != nil {
		return err
	}

	// Map [main] to config.
	err = iniFile.Section("main").MapTo(config)
	if err != nil {
		return err
	}

	// Map/Load Device Configs
	for _, section := range iniFile.Sections() {
		if section.Name() == "main" || section.Name() == "DEFAULT" {
			continue
		}
		if !section.HasKey("Type") {
			return fmt.Errorf("device [%v] must have a type defined", section.Name())
		}

		deviceConfig := DeviceConfig{
			FailureWarning: 5,
			SSHAuthMethod:  "Password",
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
	return nil
}
