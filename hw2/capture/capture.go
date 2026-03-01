package capture

import (
	"fmt"
	"time"

	"github.com/arseniizxc/network-hw2/utils"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

func CaptureArpPackets(interfaceName string) error {
	handle, err := pcap.OpenLive(interfaceName, 65536, true, pcap.BlockForever)
	if err != nil {
		return fmt.Errorf("ошибка открытия интерфейса: %v", err)
	}
	defer handle.Close()

	if err := handle.SetBPFFilter("arp"); err != nil {
		return fmt.Errorf("ошибка установки фильтра: %v", err)
	}

	fmt.Println("Захватываю ARP пакеты")

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	packetCount := 0

	for packet := range packetSource.Packets() {
		packetCount++

		ethLayer := packet.Layer(layers.LayerTypeEthernet)
		if ethLayer == nil {
			continue
		}
		eth := ethLayer.(*layers.Ethernet)

		arpLayer := packet.Layer(layers.LayerTypeARP)
		if arpLayer == nil {
			continue
		}
		arp := arpLayer.(*layers.ARP)

		timestamp := packet.Metadata().Timestamp.Format("15:04:05.000000")

		fmt.Printf("Пакет #%d:\n%s\n", packetCount, utils.FormatArpPacket(arp, eth, timestamp))
	}

	return nil
}

func CaptureArpPacketsForDuration(interfaceName string, duration time.Duration, packetChan chan<- gopacket.Packet) error {
	handle, err := pcap.OpenLive(interfaceName, 65536, true, pcap.BlockForever)
	if err != nil {
		return fmt.Errorf("ошибка открытия интерфейса: %v", err)
	}
	defer handle.Close()

	if err := handle.SetBPFFilter("arp"); err != nil {
		return fmt.Errorf("ошибка установки фильтра: %v", err)
	}

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	timer := time.NewTimer(duration)
	defer timer.Stop()

	for {
		select {
		case packet := <-packetSource.Packets():
			packetChan <- packet
		case <-timer.C:
			close(packetChan)
			return nil
		}
	}
}

func CaptureAllPacketsForDuration(interfaceName string, duration time.Duration, packetChan chan<- gopacket.Packet) error {
	handle, err := pcap.OpenLive(interfaceName, 65536, true, pcap.BlockForever)
	if err != nil {
		return fmt.Errorf("ошибка открытия интерфейса: %v", err)
	}
	defer handle.Close()

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	timer := time.NewTimer(duration)
	defer timer.Stop()

	for {
		select {
		case packet := <-packetSource.Packets():
			packetChan <- packet
		case <-timer.C:
			close(packetChan)
			return nil
		}
	}
}
