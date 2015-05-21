package sflow

import (
	"net"

	"github.com/PreetamJinka/sflow"
)

const ClassName = "sflow"

type Class struct {
	destinationAddress net.IP
	inbound            chan sflow.Datagram
}

func (c *Class) Name() string {
	return ClassName
}

func (c *Class) Category() string {
	return "comm"
}
