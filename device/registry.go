package device

import (
	"errors"
	"net"
	"sync"

	"github.com/PreetamJinka/cistern/device/class"
	"github.com/PreetamJinka/cistern/message"
)

var ErrAddressAlreadyRegistered = errors.New("device: address already registered")

type mapIP [16]byte

type Registry struct {
	sync.Mutex

	devices map[mapIP]*Device
}

func NewRegistry() *Registry {
	return &Registry{
		Mutex:   sync.Mutex{},
		devices: map[mapIP]*Device{},
	}
}

func (r *Registry) RegisterDevice(hostname string, address net.IP) (*Device, error) {
	key := toMapIP(address)

	if _, present := r.devices[key]; present {
		return nil, ErrAddressAlreadyRegistered
	}

	d := &Device{
		hostname: hostname,
		address:  address,
		classes:  map[string]class.Class{},
		messages: make(chan *message.Message),
	}

	go d.processMessages()

	r.devices[key] = d

	return d, nil
}

func (r *Registry) Lookup(address net.IP) *Device {
	return r.devices[toMapIP(address)]
}

func toMapIP(ip net.IP) mapIP {
	mIP := mapIP{}
	copy(mIP[:], ip.To16())
	return mIP
}
