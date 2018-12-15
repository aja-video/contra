package collectors

import (
	"contra/src/configuration"
	"contra/src/utils"
	"fmt"
	"github.com/google/goexpect"
)

// Collector interface keeps things together for collection.
type Collector interface {
	BuildBatcher() ([]expect.Batcher, error)
	ParseResult(string) (string, error)
}

// CollectorSpecial is special.
type CollectorSpecial interface {
	ModifySSHConfig(config *utils.SSHConfig)
}

// MakeCollector will generate the appropriate collector based on the
// type string passed in by the configuration.
func MakeCollector(d configuration.DeviceConfig) (Collector, error) {
	switch d.Type {
	case "cisco_csb":
		return makeCiscoCsb(d), nil
	case "procurve":
		return makeProcurve(d), nil
	case "comware":
		return makeComware(d), nil
	case "pfsense":
		return makePfsense(d), nil
	case "vyatta":
		return makeVyatta(d), nil
	case "arista":
		return makeArista(d), nil
	default:
		return nil, fmt.Errorf("unrecognized collector type: %v", d.Type)
	}
}
