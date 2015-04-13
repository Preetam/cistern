package decode

import (
	"github.com/PreetamJinka/sflow"

	"bytes"
	"log"
)

type SFlowDecoder struct {
	inbound  <-chan []byte
	outbound chan sflow.Datagram
}

func NewSFlowDecoder(inbound <-chan []byte, bufferLength ...int) *SFlowDecoder {
	bufLen := 0

	if len(bufferLength) > 0 {
		bufLen = bufferLength[0]
	}

	return &SFlowDecoder{
		inbound:  inbound,
		outbound: make(chan sflow.Datagram, bufLen),
	}
}

func (d *SFlowDecoder) Outbound() chan sflow.Datagram {
	return d.outbound
}

func (d *SFlowDecoder) Run() {
	decoder := sflow.NewDecoder(nil)
	go func() {
		for buf := range d.inbound {
			r := bytes.NewReader(buf)

			decoder.Use(r)

			dgram, err := decoder.Decode()
			if err == nil {
				d.outbound <- *dgram
			} else {
				log.Println(err)
			}
		}
	}()
}
