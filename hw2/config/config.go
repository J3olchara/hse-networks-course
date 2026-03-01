package config

import (
	"fmt"
	"net"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DeviceIP         net.IP
	DeviceMAC        net.HardwareAddr
	RouterIP         net.IP
	NetworkInterface string
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("ошибка загрузки .env файла: %v", err)
	}

	deviceIPStr := os.Getenv("DEVICE_IP")
	if deviceIPStr == "" {
		return nil, fmt.Errorf("DEVICE_IP не указан в .env")
	}
	deviceIP := net.ParseIP(deviceIPStr)
	if deviceIP == nil {
		return nil, fmt.Errorf("некорректный DEVICE_IP: %s", deviceIPStr)
	}

	deviceMACStr := os.Getenv("DEVICE_MAC")
	if deviceMACStr == "" {
		return nil, fmt.Errorf("DEVICE_MAC не указан в .env")
	}
	deviceMAC, err := net.ParseMAC(deviceMACStr)
	if err != nil {
		return nil, fmt.Errorf("некорректный DEVICE_MAC: %v", err)
	}

	routerIPStr := os.Getenv("ROUTER_IP")
	if routerIPStr == "" {
		return nil, fmt.Errorf("ROUTER_IP не указан в .env")
	}
	routerIP := net.ParseIP(routerIPStr)
	if routerIP == nil {
		return nil, fmt.Errorf("некорректный ROUTER_IP: %s", routerIPStr)
	}

	networkInterface := os.Getenv("NETWORK_INTERFACE")
	if networkInterface == "" {
		return nil, fmt.Errorf("NETWORK_INTERFACE не указан в .env")
	}

	return &Config{
		DeviceIP:         deviceIP,
		DeviceMAC:        deviceMAC,
		RouterIP:         routerIP,
		NetworkInterface: networkInterface,
	}, nil
}
