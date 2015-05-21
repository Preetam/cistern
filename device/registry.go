package device

import (
	"errors"
	"net"
)

var ErrAddressAlreadyRegistered = errors.New("device: address already registered")

type mapIP [16]byte

type Registry struct {
	devices map[mapIP]*Device
}

func NewRegistry() *Registry {
	return &Registry{
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
	}

	return d, nil
}

func toMapIP(ip net.IP) mapIP {
	mIP := mapIP{}
	copy(mIP[:], ip.To16())
	return mIP
}
