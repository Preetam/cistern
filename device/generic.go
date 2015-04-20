package device

import (
	"net"
	"sync"

	"github.com/PreetamJinka/snmp"

	"github.com/PreetamJinka/cistern/state/series"
)

// Generic represents a generic device.
type Generic struct {
	hostname string
	ip       net.IP

	snmpSession *snmp.Session

	Inbound  chan sflow.Datagram
	outbound chan series.Observation
}

func NewGenericDevice() *Generic {
	return &Generic{}
}
