package class

import (
	"github.com/PreetamJinka/cistern/device/class/comm/sflow"
	"github.com/PreetamJinka/cistern/message"
)

type Class interface {
	Name() string
	Category() string
	OutboundMessages() chan *message.Message
}

var _ Class = &sflow.Class{}
