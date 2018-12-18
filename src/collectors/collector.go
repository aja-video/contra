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

// CollectorSpecial is special.
type CollectorSpecial interface {
	ModifySSHConfig(config *utils.SSHConfig)
}
