package sflow

import (
	"net"

	"github.com/PreetamJinka/sflow"
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
