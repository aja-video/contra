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

	// TODO: go-ini *can* write ini files. But, not sure if it can strategically write a single value.
	// We may not want to write the entire file since some config values may be loaded from command
	// line, which we do not want to incidentally persist.
	//
	// I think we should keep the config key the same, so a user can simply delete the hashed value,
	// put a legit key, then run it. Without having to change a key from like EmailPassEncrypted to EmailPass.
	// Therefore in order to avoid double encrypting, I'm thinking of prefixing encrypted passwords.
	// Perhaps prefix with something like:
	// SMTPPass = `~~~enc~~~WkXAVkH-bPZAAKjj5T4hy_kINJRAZckoxEpjiXS_ZhI=`
	// The risk of someone's actual password being ~~~enc~~~trololo or something is pretty low, so we can
	// just check if the password has ~~~enc~~~ at the beginning, strip it and decrypt.
	//
	//log.Println(config)
	//os.Exit(1)
	// 	encryptedpass := configuration.EncryptConfig(key, pass)
	//	pass := configuration.DecryptConfig(key, encryptedpass)

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
