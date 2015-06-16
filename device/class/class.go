package class

import (
	"github.com/Preetam/cistern/device/class/comm/sflow"
	"github.com/Preetam/cistern/message"
)

type Class interface {
	Name() string
	Category() string
	OutboundMessages() chan *message.Message
}

var _ Class = &sflow.Class{}
