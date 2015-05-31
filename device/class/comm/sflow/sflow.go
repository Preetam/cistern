package sflow

import (
	"log"
	"net"

	"github.com/PreetamJinka/cistern/message"
	"github.com/PreetamJinka/sflow"
)

const ClassName = "sflow"

type Class struct {
	sourceAddress net.IP
	inbound       chan *sflow.Datagram
	outbound      chan *message.Message
}

func NewClass(sourceAddress net.IP, inbound chan *sflow.Datagram, outbound chan *message.Message) *Class {
	c := &Class{
		sourceAddress: sourceAddress,
		inbound:       inbound,
		outbound:      outbound,
	}

	go c.generateMessages()

	return c
}

func (c *Class) Name() string {
	return ClassName
}

func (c *Class) Category() string {
	return "comm"
}

func (c *Class) OutboundMessages() chan *message.Message {
	return c.outbound
}

func (c *Class) generateMessages() {
	for dgram := range c.inbound {
		log.Println("got datagram:", dgram)

		for _, sample := range dgram.Samples {
			for _, record := range sample.GetRecords() {
				if record.RecordType() == sflow.TypeHostCPUCountersRecord {
					c.outbound <- &message.Message{
						Class:   "host-counters",
						Type:    "CPU",
						Content: record,
					}
				}
			}
		}
	}
}
