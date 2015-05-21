package snmp

import (
	"github.com/PreetamJinka/snmp"
)

const ClassName = "snmp"

type Class struct {
	session *snmp.Session
}

func (c *Class) Name() string {
	return ClassName
}

func (c *Class) Category() string {
	return "comm"
}
