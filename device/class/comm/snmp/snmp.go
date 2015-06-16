package snmp

import (
	"github.com/Preetam/cistern/message"
	"github.com/Preetam/snmp"
)

const ClassName = "snmp"

type Class struct {
	session  *snmp.Session
	outbound chan *message.Message
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
