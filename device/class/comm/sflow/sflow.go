package sflow

import (
	"net"

	"github.com/PreetamJinka/sflow"
)

const ClassName = "sflow"

type Class struct {
	sourceAddress net.IP
	inbound       chan *sflow.Datagram
}

func NewClass(sourceAddress net.IP, inbound chan *sflow.Datagram) *Class {
	return &Class{
		sourceAddress: sourceAddress,
		inbound:       inbound,
	}
}

func (c *Class) Name() string {
	return ClassName
}

func (c *Class) Category() string {
	return "comm"
}
