package utils

import (
	"github.com/go-ini/ini"
	"log"
)

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
