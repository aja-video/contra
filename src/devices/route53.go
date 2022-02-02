package devices

import (
	"github.com/aja-video/contra/src/configuration"
	"github.com/aja-video/contra/src/utils"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/google/goexpect"
)

// DeviceRoute53 logic container for device.
type DeviceRoute53 struct {
	configuration.DeviceConfig
}

// SetDeviceConfig since it is unclear how to assign DeviceConfig via reflect.New
func (p *DeviceRoute53) SetDeviceConfig(deviceConfig configuration.DeviceConfig) {
	p.DeviceConfig = deviceConfig
}

// BuildBatcher for Route53
func (p *DeviceRoute53) BuildBatcher() ([]expect.Batcher, error) {
	// no expect neccesary so return a stub
	// TODO: Might be nice to rework logic to make this an optional method
	return utils.SimpleBatcher([][]string{})
}

// ParseResult for Route53
func (p *DeviceRoute53) ParseResult(result string) (string, error) {
	// build Route53 session
	mySession := session.Must(session.NewSession())
	svc := route53.New(mySession)
	// gather record sets
	input := &route53.ListResourceRecordSetsInput{
		HostedZoneId: &p.ZoneID,
	}
	records, err := svc.ListResourceRecordSets(input)
	return records.GoString(), err
}
