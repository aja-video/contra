package configuration

import (
	"fmt"
	"github.com/go-ini/ini"
	"reflect"
)

// configTester validates Contras configuration
func configTester(configFile string) error {
	var err error

	// load config file
	iniFile, err := ini.Load(configFile)
	if err != nil {
		return err
	}

	// Pull valid values for main configuration
	mainConfig := Config{}
	mainValues := reflect.ValueOf(&mainConfig).Elem()
	mainMap := buildSectionMap(mainValues, make(map[string]int8))

	// Pull valid values for device configuration
	deviceConfig := DeviceConfig{}
	deviceValues := reflect.ValueOf(&deviceConfig).Elem()
	deviceMap := buildSectionMap(deviceValues, make(map[string]int8))

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
func sectionTester(testMap map[string]int8, section *ini.Section) error {
	for _, i := range section.Keys() {
		if testMap[i.Name()] == 1 {
			if i.String() != "true" && i.String() != "false" {
				return fmt.Errorf(`invalid boolean value "%s" in section [%s] key "%s"`, i.Value(),
					section.Name(), i.Name())
			}
		} else if testMap[i.Name()] != 2 {
			return fmt.Errorf(`invalid key "%s" in section [%s]`, i.Name(), section.Name())
		}
	}
	return nil
}

// initialize map for config section
func buildSectionMap(value reflect.Value, testMap map[string]int8) map[string]int8 {
	for i := 0; i < value.NumField(); i++ {
		key := value.Type().Field(i).Name
		// set value to 1 for boolean so we can identify them
		if value.Type().Field(i).Type.String() == "bool" {
			testMap[key] = 1
		} else {
			testMap[key] = 2
		}
	}
	return testMap
}
