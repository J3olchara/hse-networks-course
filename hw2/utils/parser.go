package utils

import (
	"bytes"
	"fmt"
	"net"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

func IsGratuitousArp(arp *layers.ARP) bool {
	if arp.Operation != layers.ARPRequest {
		return false
	}
	
	senderIP := net.IP(arp.SourceProtAddress)
	targetIP := net.IP(arp.DstProtAddress)
	
	if !senderIP.Equal(targetIP) {
		return false
	}
	
	targetMAC := net.HardwareAddr(arp.DstHwAddress)
	zeroMAC := net.HardwareAddr{0, 0, 0, 0, 0, 0}
	broadcastMAC := net.HardwareAddr{0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
	
	return bytes.Equal(targetMAC, zeroMAC) || bytes.Equal(targetMAC, broadcastMAC)
}

func IsBroadcast(eth *layers.Ethernet) bool {
	broadcastMAC := net.HardwareAddr{0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
	return bytes.Equal(eth.DstMAC, broadcastMAC)
}

func FormatArpPacket(arp *layers.ARP, eth *layers.Ethernet, timestamp string) string {
	var operation string
	switch arp.Operation {
	case layers.ARPRequest:
		operation = "REQUEST"
	case layers.ARPReply:
		operation = "REPLY"
	default:
		operation = fmt.Sprintf("UNKNOWN (%d)", arp.Operation)
	}

	senderMAC := net.HardwareAddr(arp.SourceHwAddress)
	senderIP := net.IP(arp.SourceProtAddress)
	targetMAC := net.HardwareAddr(arp.DstHwAddress)
	targetIP := net.IP(arp.DstProtAddress)

	return fmt.Sprintf(`Время: %s
Ethernet:
  Source MAC: %s
  Destination MAC: %s
ARP:
  Operation: %s
  Sender MAC: %s
  Sender IP: %s
  Target MAC: %s
  Target IP: %s
`,
		timestamp,
		eth.SrcMAC,
		eth.DstMAC,
		operation,
		senderMAC,
		senderIP,
		targetMAC,
		targetIP,
	)
}

func CalculateFrameSize(packet gopacket.Packet) int {
	captureLength := packet.Metadata().CaptureLength
	
	if captureLength < 64 {
		return 64
	}
	
	return captureLength
}

func IsDeviceToRouter(eth *layers.Ethernet, deviceMAC, routerMAC net.HardwareAddr) bool {
	return bytes.Equal(eth.SrcMAC, deviceMAC) && bytes.Equal(eth.DstMAC, routerMAC)
}

func IsRouterToDevice(eth *layers.Ethernet, deviceMAC, routerMAC net.HardwareAddr) bool {
	return bytes.Equal(eth.SrcMAC, routerMAC) && bytes.Equal(eth.DstMAC, deviceMAC)
}
