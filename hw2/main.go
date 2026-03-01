package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/arseniizxc/network-hw2/capture"
	"github.com/arseniizxc/network-hw2/config"
	"github.com/arseniizxc/network-hw2/sender"
	"github.com/arseniizxc/network-hw2/statistics"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf(err.Error())
		os.Exit(1)
	}

	fmt.Printf("IP устройства: %s\n", cfg.DeviceIP)
	fmt.Printf("MAC устройства: %s\n", cfg.DeviceMAC)
	fmt.Printf("IP роутера: %s\n", cfg.RouterIP)
	fmt.Printf("Интерфейс: %s\n\n", cfg.NetworkInterface)

	fmt.Println("Программу нужно запускать с sudo")

	reader := bufio.NewReader(os.Stdin)

	for {
		printMenu()

		fmt.Print("Выберите команду: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Ошибка чтения ввода: %v\n", err)
			continue
		}

		input = strings.TrimSpace(input)

		switch input {
		case "1":
			fmt.Println("Захват ARP пакетов")
			if err := capture.CaptureArpPackets(cfg.NetworkInterface); err != nil {
				fmt.Printf("Ошибка: %v\n", err)
			}

		case "2":
			fmt.Println("Ищу MAC адрес роутера")
			routerMAC, err := sender.FindRouterMAC(
				cfg.NetworkInterface,
				cfg.DeviceIP,
				cfg.DeviceMAC,
				cfg.RouterIP,
			)
			if err != nil {
				fmt.Printf("Ошибка: %v\n", err)
			} else {
				fmt.Printf("MAC адрес роутера %s: %s\n", cfg.RouterIP, routerMAC)
			}

		case "3":
			fmt.Println("Сбор статистики")
			fmt.Print("Введите время в секундах: ")

			durationInput, err := reader.ReadString('\n')
			if err != nil {
				fmt.Printf("Ошибка чтения ввода: %v\n", err)
				continue
			}

			durationInput = strings.TrimSpace(durationInput)
			seconds, err := strconv.Atoi(durationInput)
			if err != nil || seconds <= 0 {
				fmt.Println("Нужно ввести положительное число")
				continue
			}

			routerMAC, err := sender.FindRouterMAC(
				cfg.NetworkInterface,
				cfg.DeviceIP,
				cfg.DeviceMAC,
				cfg.RouterIP,
			)
			if err != nil {
				fmt.Printf("Не удалось определить MAC роутера: %v\n", err)
				fmt.Println("Статистика будет собрана без учета трафика с роутером")
				routerMAC = nil
			}

			duration := time.Duration(seconds) * time.Second
			if err := statistics.CollectStatistics(cfg.NetworkInterface, duration, cfg.DeviceMAC, routerMAC); err != nil {
				fmt.Printf("Ошибка: %v\n", err)
			}

		case "4", "exit", "quit", "q":
			fmt.Println("Выход")
			os.Exit(0)

		default:
			fmt.Printf("Неизвестная команда: %s\n", input)
		}

		fmt.Println()
	}
}

func printMenu() {
	fmt.Println("1 - Захват ARP пакетов")
	fmt.Println("2 - Найти MAC адрес роутера")
	fmt.Println("3 - Собрать статистику")
	fmt.Println("4 - Выход")
}
