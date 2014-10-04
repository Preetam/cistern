package main

import (
	"github.com/PreetamJinka/protodecode"
	"github.com/PreetamJinka/sflow-go"

	"fmt"
)

type RawPacketProcessor struct {
	reg      *HostRegistry
	inbound  chan Message
	outbound chan Message
}

func NewRawPacketProcessor(reg *HostRegistry) *RawPacketProcessor {
	return &RawPacketProcessor{
		reg:      reg,
		outbound: make(chan Message, 100),
	}
}

func (p *RawPacketProcessor) SetInbound(inbound chan Message) {
	p.inbound = inbound
}

func (p *RawPacketProcessor) Outbound() chan Message {
	return p.outbound
}

func (p *RawPacketProcessor) Process() {
	for message := range p.inbound {
		record := message.Record

		switch record.RecordType() {
		case sflow.TypeRawPacketFlow:
			r := record.(sflow.RawPacketFlowRecord)

			b := make([]byte, len(r.Header))
			copy(b, r.Header)

			ethernetPacket := protodecode.DecodeEthernet(b)
			switch ethernetPacket.EtherType {
			case 0x0800:
				ipv4Packet := protodecode.DecodeIPv4(ethernetPacket.Payload)

				switch ipv4Packet.Protocol {
				case 6: // TCP
					tcpPacket := protodecode.DecodeTCP(ipv4Packet.Payload)
					fmt.Printf("[TCP] %s:%d => %s:%d [%d bytes]\n", ipv4Packet.Source, tcpPacket.SourcePort,
						ipv4Packet.Destination, tcpPacket.DestinationPort,
						ipv4Packet.Length)

					if tcpPacket.HasSYN() && !tcpPacket.HasACK() {
						fmt.Println("  connection opening")
					}

					if tcpPacket.HasFIN() && tcpPacket.HasACK() {
						fmt.Println("  connection closing")
					}

				case 17: // UDP
					udpPacket := protodecode.DecodeUDP(ipv4Packet.Payload)
					fmt.Printf("[UDP] %s:%d => %s:%d [%d bytes]\n", ipv4Packet.Source, udpPacket.SourcePort,
						ipv4Packet.Destination, udpPacket.DestinationPort,
						ipv4Packet.Length)

				default:
					fmt.Printf("[???] %s => %s [%d bytes]\n", ipv4Packet.Source, ipv4Packet.Destination,
						ipv4Packet.Length)
				}

			case 0x86dd:
				ipv6Packet := protodecode.DecodeIPv6(ethernetPacket.Payload)

				switch ipv6Packet.NextHeader {
				case 6: // TCP
					tcpPacket := protodecode.DecodeTCP(ipv6Packet.Payload)
					fmt.Printf("[TCP] %s:%d => %s:%d [%d bytes]\n", ipv6Packet.SourceAddress, tcpPacket.SourcePort,
						ipv6Packet.DestinationAddress, tcpPacket.DestinationPort,
						ipv6Packet.PayloadLength)

					if tcpPacket.HasSYN() && !tcpPacket.HasACK() {
						fmt.Println("  connection opening")
					}

					if tcpPacket.HasFIN() && tcpPacket.HasACK() {
						fmt.Println("  connection closing")
					}

				case 17: // UDP
					udpPacket := protodecode.DecodeUDP(ipv6Packet.Payload)
					fmt.Printf("[UDP] %s:%d => %s:%d [%d bytes]\n", ipv6Packet.SourceAddress, udpPacket.SourcePort,
						ipv6Packet.DestinationAddress, udpPacket.DestinationPort,
						ipv6Packet.PayloadLength)

				default:
					fmt.Printf("[???] %s => %s [%d bytes]\n", ipv6Packet.SourceAddress, ipv6Packet.DestinationAddress,
						ipv6Packet.PayloadLength)
				}
			}

		default:
			select {
			case p.outbound <- message:
			default:
			}
		}
	}
}