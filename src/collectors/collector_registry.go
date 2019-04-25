package collectors

import "github.com/aja-video/contra/src/devices"

// Mandatory that new collector definitions be added to this array.
var deviceMap = map[string]interface{}{
	"arista":    devices.DeviceArista{},
	"cisco_csb": devices.DeviceCiscoCsb{},
	"comware":   devices.DeviceComware{},
	"mikrotik":  devices.DeviceMikrotik{},
	"pfsense":   devices.DevicePfsense{},
	"procurve":  devices.DeviceProcurve{},
	"vyatta":    devices.DeviceVyatta{},
}
