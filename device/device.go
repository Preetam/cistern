package device

import (
	"errors"

	"github.com/PreetamJinka/cistern/device/class"
)

import (
	"net"
	"sync"
)

var (
	ErrClassNameRegistered = errors.New("device: class name already registered")
)

type Device struct {
	sync.Mutex

	hostname string
	address  net.IP

	classes map[string]class.Class
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
