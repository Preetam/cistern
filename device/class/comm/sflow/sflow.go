package sflow

import (
	"github.com/PreetamJinka/sflow"

	"net"
)

type Class struct {
	destinationAddress net.IP
	inbound            chan sflow.Datagram
}

func (c *Class) Name() string {
	return "sflow"
}

func (c *Class) Category() string {
	return "comm"
}
