package utils

import (
	"net"
	"time"
)

type NetworkConfig struct {
	DeviceIP         net.IP
	DeviceMAC        net.HardwareAddr
	RouterIP         net.IP
	NetworkInterface string
}

type ArpRequest struct {
	SenderIP   net.IP
	TargetIP   net.IP
	Timestamp  time.Time
}

type Statistics struct {
	TotalEthernetFrames   int
	TotalArpPackets       int
	UniqueMACAddresses    map[string]bool
	BroadcastFrames       int
	BroadcastArpFrames    int
	GratuitousArpRequests int
	ArpRequestResponsePairs int
	DataVolumeToRouter    int64
	DataVolumeFromRouter  int64
}

func NewStatistics() *Statistics {
	return &Statistics{
		UniqueMACAddresses: make(map[string]bool),
	}
}
