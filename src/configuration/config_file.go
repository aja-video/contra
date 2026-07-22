package configuration

import (
	"fmt"
	"github.com/go-ini/ini"
	"time"
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
			FailureWarning:         5,
			FailureBackoffCount:    5,
			FailureBackoffInterval: 24 * time.Hour,
			SSHAuthMethod:          "Password",
			SSHTimeout:             10 * time.Second,
			AllowInsecureSSH:       true,
		}

		section.MapTo(&deviceConfig)
		// Copy the section name into the device config for reference.
		deviceConfig.Name = section.Name()
		config.Devices = append(config.Devices, deviceConfig)
	}
	return nil
}
