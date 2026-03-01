# ARP Protocol Tool

Приложение для работы с протоколом ARP с использованием библиотеки PCAP.

## Описание

Это консольное приложение позволяет:
1. Захватывать и интерпретировать ARP пакеты в PROMISCUOUS режиме
2. Определять MAC адрес роутера через ARP запрос
3. Собирать статистику по ARP трафику за заданное время

## Требования

### Системные требования

- **ОС**: Linux
- **Права**: sudo (для работы в PROMISCUOUS режиме)
- **Библиотеки**: libpcap-dev

### Установка зависимостей

```bash
# Ubuntu/Debian
sudo apt-get update
sudo apt-get install libpcap-dev

# Fedora/RHEL
sudo dnf install libpcap-devel

# Arch Linux
sudo pacman -S libpcap
```

### Go зависимости

- Go 1.21 или выше
- github.com/google/gopacket v1.1.19
- github.com/joho/godotenv v1.5.1

## Установка и настройка

### 1. Клонирование/распаковка проекта

```bash
cd hw2
```

### 2. Настройка конфигурации

Скопируйте файл `.env.example` в `.env` и заполните своими значениями:

```bash
cp .env.example .env
nano .env
```

#### Как узнать свои сетевые параметры:

**IP адрес и MAC адрес устройства:**
```bash
ip addr show
```
Найдите ваш активный интерфейс (обычно eth0, wlan0, enp0s3), в выводе будут:
- `inet 192.168.1.X/24` - это ваш IP
- `link/ether aa:bb:cc:dd:ee:ff` - это ваш MAC

**IP адрес роутера:**
```bash
ip route | grep default
```
Вывод: `default via 192.168.1.1 dev eth0` - 192.168.1.1 это IP роутера

**Название интерфейса:**
Обычно `eth0`, `wlan0`, `enp0s3` - смотрите в выводе `ip addr show`

Пример `.env`:
```env
DEVICE_IP=192.168.1.100
DEVICE_MAC=aa:bb:cc:dd:ee:ff
ROUTER_IP=192.168.1.1
NETWORK_INTERFACE=eth0
```

### 3. Установка Go зависимостей

```bash
go mod download
```

### 4. Компиляция

```bash
go build -o arp-tool
```

## Запуск

**ВАЖНО:** Приложение требует прав root для работы с PCAP в PROMISCUOUS режиме:

```bash
sudo ./arp-tool
```

## Использование

После запуска программа выведет меню с доступными командами:

### Команда 1: Захват ARP пакетов

Захватывает все ARP пакеты в сети и выводит их интерпретацию на консоль.

**Что отображается:**
- Ethernet заголовок (Source/Destination MAC, EtherType)
- ARP заголовок (все поля)
- Timestamp пакета
- Тип операции (REQUEST/REPLY)

**Остановка:** Нажмите `Ctrl+C`

### Команда 2: Найти MAC адрес роутера

Отправляет ARP запрос на IP роутера (из `.env`) и выводит полученный MAC адрес.

**Вывод:**
- IP роутера
- Найденный MAC адрес
- Время получения ответа

### Команда 3: Собрать статистику

Собирает статистику по трафику за заданное пользователем время.

**Собираемые показатели:**
1. Количество Ethernet фреймов
2. Количество ARP пакетов
3. Количество уникальных MAC адресов (с выводом списка)
4. Количество широковещательных Ethernet фреймов
5. Количество широковещательных Ethernet фреймов с ARP
6. Количество Gratuitous ARP Requests
7. Количество пар ARP Request/Response
8. Объем данных между устройством и роутером (в байтах)

**Рекомендации:**
- Запустите сбор на 30-60 секунд
- Во время сбора генерируйте трафик: открывайте сайты, пингуйте устройства
- Параллельно запустите WireShark для верификации

### Команда 4: Выход

Завершает работу программы.

## Проверка работы с WireShark

### Установка WireShark

```bash
sudo apt-get install wireshark
sudo usermod -a -G wireshark $USER
```

### Запуск

```bash
sudo wireshark
```

### Полезные фильтры

- `arp` - все ARP пакеты
- `arp.opcode == 1` - только ARP запросы
- `arp.opcode == 2` - только ARP ответы
- `eth.dst == ff:ff:ff:ff:ff:ff` - широковещательные фреймы
- `arp.src.proto_ipv4 == arp.dst.proto_ipv4` - Gratuitous ARP
- `eth.addr == XX:XX:XX:XX:XX:XX` - пакеты с конкретным MAC

## Структура проекта

```
hw2/
├── main.go              - Консольное меню и точка входа
├── config/
│   └── config.go        - Загрузка конфигурации из .env
├── capture/
│   └── capture.go       - Захват ARP пакетов
├── sender/
│   └── sender.go        - Отправка ARP запросов
├── statistics/
│   └── statistics.go    - Сбор и вывод статистики
└── utils/
    ├── types.go         - Структуры данных
    └── parser.go        - Утилиты для анализа пакетов
```

См. `FUNCTIONS.md` для детального списка всех функций и мест их реализации.

## Решение проблем

### "Permission denied" при запуске

Приложение требует прав root:
```bash
sudo ./arp-tool
```

### "No such device"

Проверьте название интерфейса в `.env`:
```bash
ip addr show
```

### "pcap.h: No such file or directory" при компиляции

Установите libpcap-dev:
```bash
sudo apt-get install libpcap-dev
```

### ARP пакеты не захватываются

1. Проверьте, что интерфейс активен: `ip link show`
2. Убедитесь, что запускаете с sudo
3. Попробуйте сгенерировать ARP трафик: `ping [другой IP в сети]`

## Авторы

Выполнено в рамках HW2 по курсу "Компьютерные сети"

## Лицензия

Учебный проект
