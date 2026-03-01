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

	fmt.Printf("device ip:  %s\n", cfg.DeviceIP)
	fmt.Printf("device mac: %s\n", cfg.DeviceMAC)
	fmt.Printf("router ip: %s\n", cfg.RouterIP)
	fmt.Printf("network interface: %s\n\n", cfg.NetworkInterface)

	fmt.Println("обязательно запустите программу с sudo")

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
			fmt.Println("\n--- Команда 1: Захват ARP пакетов ---")
			if err := capture.CaptureArpPackets(cfg.NetworkInterface); err != nil {
				fmt.Printf("ОШИБКА: %v\n", err)
			}

		case "2":
			fmt.Println("\n--- Команда 2: Найти MAC адрес роутера ---")
			routerMAC, err := sender.FindRouterMAC(
				cfg.NetworkInterface,
				cfg.DeviceIP,
				cfg.DeviceMAC,
				cfg.RouterIP,
			)
			if err != nil {
				fmt.Printf("ОШИБКА: %v\n", err)
			} else {
				fmt.Printf("MAC адрес роутера %s: %s\n", cfg.RouterIP, routerMAC)
			}

		case "3":
			fmt.Println("\n--- Команда 3: Собрать статистику ---")
			fmt.Print("Введите время сбора статистики (в секундах): ")

			durationInput, err := reader.ReadString('\n')
			if err != nil {
				fmt.Printf("Ошибка чтения ввода: %v\n", err)
				continue
			}

			durationInput = strings.TrimSpace(durationInput)
			seconds, err := strconv.Atoi(durationInput)
			if err != nil || seconds <= 0 {
				fmt.Println("ОШИБКА: Введите корректное положительное число секунд")
				continue
			}

			routerMAC, err := sender.FindRouterMAC(
				cfg.NetworkInterface,
				cfg.DeviceIP,
				cfg.DeviceMAC,
				cfg.RouterIP,
			)
			if err != nil {
				fmt.Printf("ОШИБКА определения MAC роутера: %v\n", err)
				fmt.Println("Статистика будет собрана без учета трафика с роутером")
				routerMAC = nil
			}

			duration := time.Duration(seconds) * time.Second
			if err := statistics.CollectStatistics(cfg.NetworkInterface, duration, cfg.DeviceMAC, routerMAC); err != nil {
				fmt.Printf("ОШИБКА: %v\n", err)
			}

		case "4", "exit", "quit", "q":
			fmt.Println("\nВыход из программы. До свидания!")
			os.Exit(0)

		default:
			fmt.Printf("Неизвестная команда: %s\n", input)
		}

		fmt.Println()
	}
}

func printMenu() {
	fmt.Println("1 - Захват всех ARP пакетов (PROMISCUOUS режим)")
	fmt.Println("    • Интерпретация заголовков ARP")
	fmt.Println("    • Вывод всех полей пакета")
	fmt.Println("    • Нажмите Ctrl+C для остановки")
	fmt.Println()
	fmt.Println("2 - Найти MAC адрес роутера")
	fmt.Println("    • Отправка ARP запроса")
	fmt.Println("    • Получение и вывод MAC адреса")
	fmt.Println()
	fmt.Println("3 - Собрать статистику за заданное время")
	fmt.Println("    • Количество Ethernet фреймов и ARP пакетов")
	fmt.Println("    • Уникальные MAC адреса")
	fmt.Println("    • Широковещательные сообщения")
	fmt.Println("    • Gratuitous ARP Requests")
	fmt.Println("    • Пары ARP Request/Response")
	fmt.Println("    • Объем данных между устройством и роутером")
	fmt.Println()
	fmt.Println("4 - Выход")
}
