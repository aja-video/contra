package configuration

import (
	"fmt"
	"github.com/go-ini/ini"
	"reflect"
)

// configTester validates Contras configuration
// TODO: Sanity check values?
func configTester(configFile string) error {
	var mainMap = make(map[string]bool)
	var deviceMap = make(map[string]bool)
	var err error

	// load config file
	iniFile, err := ini.Load(configFile)
	if err != nil {
		return err
	}

	// Pull valid values for main configuration
	mainConfig := Config{}
	mainValues := reflect.ValueOf(&mainConfig).Elem()
	for i := 0; i < mainValues.NumField(); i++ {
		key := mainValues.Type().Field(i).Name
		mainMap[key] = true
	}

	// Pull valid values for device configuration
	deviceConfig := DeviceConfig{}
	deviceValues := reflect.ValueOf(&deviceConfig).Elem()
	for i := 0; i < deviceValues.NumField(); i++ {
		key := deviceValues.Type().Field(i).Name
		deviceMap[key] = true
	}

	// iterate through ini sections and check them against valid config keys
	for _, section := range iniFile.Sections() {
		if section.Name() == "main" {
			err = sectionTester(mainMap, section)
		} else {
			err = sectionTester(deviceMap, section)
		}
		if err != nil {
			return err
		}

	}

	return nil
}

// sectionTester tests a ini config section against a map of valid keys
func sectionTester(testMap map[string]bool, section *ini.Section) error {
	for _, i := range section.KeyStrings() {
		if !testMap[i] {
			return fmt.Errorf(`invalid key "%s" in section [%s]`, i, section.Name())
		}
	}
	return nil
}
