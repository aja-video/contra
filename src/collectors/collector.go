package collectors

import (
	"github.com/aja-video/contra/src/configuration"
	"github.com/aja-video/contra/src/utils"
	"github.com/google/goexpect"
)

// Collector interface keeps things together for collection.
type Collector interface {
	SetDeviceConfig(d configuration.DeviceConfig)
	BuildBatcher() ([]expect.Batcher, error)
	ParseResult(string) (string, error)
}

// CollectorSpecialSSH is special.
type CollectorSpecialSSH interface {
	ModifySSHConfig(config *utils.SSHConfig)
}

// CollectorSpecialTerminal is also special
type CollectorSpecialTerminal interface {
	ModifyUsername(config *utils.SSHConfig)
}
