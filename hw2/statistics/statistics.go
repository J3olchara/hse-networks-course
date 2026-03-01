package statistics

import (
	"fmt"
	"net"
	"sort"
	"time"

	"github.com/arseniizxc/network-hw2/capture"
	"github.com/arseniizxc/network-hw2/utils"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

type requestInfo struct {
	senderIP  net.IP
	targetIP  net.IP
	timestamp time.Time
}

func CollectStatistics(interfaceName string, duration time.Duration, deviceMAC, routerMAC net.HardwareAddr) error {
	fmt.Printf("Собираю статистику в течение %d секунд\n", int(duration.Seconds()))
	fmt.Println("Для генерации трафика можно пинговать устройства или открыть браузер")

	stats := utils.NewStatistics()
	pendingRequests := make(map[string]requestInfo)
	
	packetChan := make(chan gopacket.Packet, 1000)
	
	go func() {
		err := capture.CaptureAllPacketsForDuration(interfaceName, duration, packetChan)
		if err != nil {
			fmt.Printf("Ошибка захвата: %v\n", err)
		}
	}()

	startTime := time.Now()
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case packet, ok := <-packetChan:
			if !ok {
				printStatistics(stats, deviceMAC, routerMAC)
				return nil
			}
			
			processPacket(packet, stats, pendingRequests, deviceMAC, routerMAC)

		case <-ticker.C:
			elapsed := time.Since(startTime).Seconds()
			fmt.Printf("Прошло %.0f секунд, собрано %d Ethernet фреймов, %d ARP пакетов\n", 
				elapsed, stats.TotalEthernetFrames, stats.TotalArpPackets)
		}
	}
}

func processPacket(packet gopacket.Packet, stats *utils.Statistics, pendingRequests map[string]requestInfo, deviceMAC, routerMAC net.HardwareAddr) {
	ethLayer := packet.Layer(layers.LayerTypeEthernet)
	if ethLayer == nil {
		return
	}
	eth := ethLayer.(*layers.Ethernet)
	
	stats.TotalEthernetFrames++
	
	stats.UniqueMACAddresses[eth.SrcMAC.String()] = true
	stats.UniqueMACAddresses[eth.DstMAC.String()] = true
	
	if utils.IsBroadcast(eth) {
		stats.BroadcastFrames++
	}
	
	if utils.IsDeviceToRouter(eth, deviceMAC, routerMAC) {
		stats.DataVolumeToRouter += int64(utils.CalculateFrameSize(packet))
	} else if utils.IsRouterToDevice(eth, deviceMAC, routerMAC) {
		stats.DataVolumeFromRouter += int64(utils.CalculateFrameSize(packet))
	}
	
	arpLayer := packet.Layer(layers.LayerTypeARP)
	if arpLayer == nil {
		return
	}
	arp := arpLayer.(*layers.ARP)
	
	stats.TotalArpPackets++
	
	if utils.IsBroadcast(eth) {
		stats.BroadcastArpFrames++
	}
	
	if utils.IsGratuitousArp(arp) {
		stats.GratuitousArpRequests++
	}
	
	if arp.Operation == layers.ARPRequest {
		senderIP := net.IP(arp.SourceProtAddress)
		targetIP := net.IP(arp.DstProtAddress)
		key := fmt.Sprintf("%s->%s", senderIP, targetIP)
		pendingRequests[key] = requestInfo{
			senderIP:  senderIP,
			targetIP:  targetIP,
			timestamp: packet.Metadata().Timestamp,
		}
	} else if arp.Operation == layers.ARPReply {
		senderIP := net.IP(arp.SourceProtAddress)
		targetIP := net.IP(arp.DstProtAddress)
		key := fmt.Sprintf("%s->%s", targetIP, senderIP)
		
		if req, exists := pendingRequests[key]; exists {
			timeDiff := packet.Metadata().Timestamp.Sub(req.timestamp)
			if timeDiff < 5*time.Second {
				stats.ArpRequestResponsePairs++
				delete(pendingRequests, key)
			}
		}
	}
	
	cleanupOldRequests(pendingRequests, packet.Metadata().Timestamp)
}

func cleanupOldRequests(pendingRequests map[string]requestInfo, currentTime time.Time) {
	for key, req := range pendingRequests {
		if currentTime.Sub(req.timestamp) > 5*time.Second {
			delete(pendingRequests, key)
		}
	}
}

func printStatistics(stats *utils.Statistics, deviceMAC, routerMAC net.HardwareAddr) {
	fmt.Println("\nРезультаты:")
	
	fmt.Printf("Всего Ethernet фреймов: %d\n", stats.TotalEthernetFrames)
	fmt.Printf("Всего ARP пакетов: %d\n", stats.TotalArpPackets)
	
	fmt.Printf("Уникальных MAC адресов: %d\n", len(stats.UniqueMACAddresses))
	macList := make([]string, 0, len(stats.UniqueMACAddresses))
	for mac := range stats.UniqueMACAddresses {
		macList = append(macList, mac)
	}
	sort.Strings(macList)
	for _, mac := range macList {
		label := ""
		if mac == deviceMAC.String() {
			label = " (ваше устройство)"
		} else if mac == routerMAC.String() {
			label = " (роутер)"
		}
		fmt.Printf("  %s%s\n", mac, label)
	}
	
	fmt.Printf("Широковещательных Ethernet фреймов: %d\n", stats.BroadcastFrames)
	fmt.Printf("Из них с протоколом ARP: %d\n", stats.BroadcastArpFrames)
	
	fmt.Printf("Gratuitous ARP запросов: %d\n", stats.GratuitousArpRequests)
	
	fmt.Printf("Пар ARP запрос-ответ: %d\n", stats.ArpRequestResponsePairs)
	
	totalData := stats.DataVolumeToRouter + stats.DataVolumeFromRouter
	fmt.Printf("Объем данных между устройством и роутером: %d байт\n", totalData)
	fmt.Printf("  Устройство -> Роутер: %d байт\n", stats.DataVolumeToRouter)
	fmt.Printf("  Роутер -> Устройство: %d байт\n", stats.DataVolumeFromRouter)
}
