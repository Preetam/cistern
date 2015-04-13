package device

import (
	"net"

	"github.com/PreetamJinka/snmp"

	"github.com/PreetamJinka/cistern/state/flows"
	"github.com/PreetamJinka/cistern/state/metrics"
)

type deviceType int

const (
	TypeUnknown deviceType = 0

	TypeNetwork deviceType = 1 << (iota - 1)
	TypeLinux
	TypeBSD
)

var (
	descOid     = snmp.MustParseOID(".1.3.6.1.2.1.1.1.0")
	hostnameOid = snmp.MustParseOID(".1.3.6.1.2.1.1.5.0")
)

// A Device is an entity that sends flows or
// makes information available via SNMP.
type Device interface {
	Hostname() string
	Desc() string
	IP() net.IP
	Type() deviceType
	Discover()
	Metrics() []metrics.MetricDefinition

	// TODO: this should really go somewhere else
	TopTalkers() *flows.TopTalkers
}
