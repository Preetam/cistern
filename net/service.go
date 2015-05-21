// Package net provides a Service, which handles all of Cistern's
// protocol and network integration.
package net

import (
	"sync"
)

type Config struct {
	SFlowAddr string `json:"sflowAddr"`
}

var DefaultConfig = Config{
	SFlowAddr: ":6343",
}

type Service struct {
	lock sync.Mutex
}

func NewService(conf Config) (*Service, error) {
	// TODO: use config

	return &Service{
		lock: sync.Mutex{},
	}, nil
}
