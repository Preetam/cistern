package net

import (
	"sync"
)

type Service struct {
	lock sync.Mutex
}

func NewService() *Service {
	return &Service{
		lock: sync.Mutex{},
	}
}
