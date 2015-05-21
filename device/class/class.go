package class

import (
	"github.com/PreetamJinka/cistern/device/class/comm/sflow"
)

type Class interface {
	Name() string
	Category() string
}

var _ Class = &sflow.Class{}
