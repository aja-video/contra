package collectors

import (
	"contra/src/utils"
	"github.com/google/goexpect"
)

// Collector interface keeps things together for collection.
type Collector interface {
	BuildBatcher() ([]expect.Batcher, error)
	ParseResult(string) (string, error)
}

// CollectorSpecial is special.
type CollectorSpecial interface {
	ModifySSHConfig(config utils.SSHConfig)
}

// CollectorDefinition write me.
type CollectorDefinition struct {
	name string
}
