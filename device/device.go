package device

import (
	"errors"
	"fmt"

	"github.com/Preetam/cistern/device/class"
)

import (
	"log"
	"net"
	"sync"

	"github.com/Preetam/cistern/message"
)

var (
	ErrClassNameRegistered = errors.New("device: class name already registered")
)

type Device struct {
	sync.Mutex

	hostname string
	address  net.IP

	classes  map[string]class.Class
	messages chan *message.Message
}

func (d *Device) RegisterClass(c class.Class) error {
	if _, present := d.classes[c.Name()]; present {
		return ErrClassNameRegistered
	}

	d.classes[c.Name()] = c
	return nil
}

func (d *Device) HasClass(classname string) bool {
	_, present := d.classes[classname]
	return present
}

func (d *Device) Messages() chan *message.Message {
	return d.messages
}

func (d *Device) processMessages() {
	for m := range d.messages {
		log.Printf("%v got message: %#+v", d, m)
	}
}

func (d *Device) String() string {
	if d.hostname == "" {
		return fmt.Sprintf("Device{%v}", d.address)
	}

	return fmt.Sprintf("Device{%s - %v}", d.hostname, d.address)
}
