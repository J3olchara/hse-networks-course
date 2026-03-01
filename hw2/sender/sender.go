package sender

import (
	"bytes"
	"fmt"
	"net"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

func FindRouterMAC(interfaceName string, deviceIP, deviceMAC net.HardwareAddr, routerIP net.IP) (net.HardwareAddr, error) {
	handle, err := pcap.OpenLive(interfaceName, 65536, true, pcap.BlockForever)
	if err != nil {
		return nil, fmt.Errorf("ошибка открытия интерфейса: %v", err)
	}
	defer handle.Close()

	startTime := time.Now()
	fmt.Printf("\nОтправка ARP запроса на IP: %s\n", routerIP)
	
	if err := sendArpRequest(handle, deviceIP, deviceMAC, routerIP); err != nil {
		return nil, fmt.Errorf("ошибка отправки ARP запроса: %v", err)
	}

	if err := handle.SetBPFFilter("arp"); err != nil {
		return nil, fmt.Errorf("ошибка установки фильтра: %v", err)
	}

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	timeout := time.After(5 * time.Second)

	for {
		select {
		case packet := <-packetSource.Packets():
			arpLayer := packet.Layer(layers.LayerTypeARP)
			if arpLayer == nil {
				continue
			}
			arp := arpLayer.(*layers.ARP)

			if arp.Operation != layers.ARPReply {
				continue
			}

			senderIP := net.IP(arp.SourceProtAddress)
			if senderIP.Equal(routerIP) {
				routerMAC := net.HardwareAddr(arp.SourceHwAddress)
				elapsed := time.Since(startTime)
				
				fmt.Printf("\nПолучен ответ:\n")
				fmt.Printf("  MAC адрес роутера: %s\n", routerMAC)
				fmt.Printf("  Время ответа: %d ms\n\n", elapsed.Milliseconds())
				
				return routerMAC, nil
			}

		case <-timeout:
			return nil, fmt.Errorf("таймаут ожидания ARP ответа от роутера")
		}
	}
}

func sendArpRequest(handle *pcap.Handle, deviceIP net.IP, deviceMAC net.HardwareAddr, targetIP net.IP) error {
	eth := layers.Ethernet{
		SrcMAC:       deviceMAC,
		DstMAC:       net.HardwareAddr{0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
		EthernetType: layers.EthernetTypeARP,
	}

	arp := layers.ARP{
		AddrType:          layers.LinkTypeEthernet,
		Protocol:          layers.EthernetTypeIPv4,
		HwAddressSize:     6,
		ProtAddressSize:   4,
		Operation:         layers.ARPRequest,
		SourceHwAddress:   []byte(deviceMAC),
		SourceProtAddress: []byte(deviceIP.To4()),
		DstHwAddress:      []byte{0, 0, 0, 0, 0, 0},
		DstProtAddress:    []byte(targetIP.To4()),
	}

	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}

	if err := gopacket.SerializeLayers(buf, opts, &eth, &arp); err != nil {
		return fmt.Errorf("ошибка сериализации пакета: %v", err)
	}

	if err := handle.WritePacketData(buf.Bytes()); err != nil {
		return fmt.Errorf("ошибка отправки пакета: %v", err)
	}

	return nil
}

func SendArpRequest(interfaceName string, deviceIP net.IP, deviceMAC net.HardwareAddr, targetIP net.IP) error {
	handle, err := pcap.OpenLive(interfaceName, 65536, false, pcap.BlockForever)
	if err != nil {
		return fmt.Errorf("ошибка открытия интерфейса: %v", err)
	}
	defer handle.Close()

	return sendArpRequest(handle, deviceIP, deviceMAC, targetIP)
}

func IsArpReplyFor(arp *layers.ARP, targetIP net.IP) bool {
	if arp.Operation != layers.ARPReply {
		return false
	}
	senderIP := net.IP(arp.SourceProtAddress)
	return bytes.Equal(senderIP, targetIP)
}
