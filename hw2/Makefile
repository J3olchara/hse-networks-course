.PHONY: build run clean install-deps help

help:
	@echo "Доступные команды:"
	@echo "  make install-deps  - Установить системные зависимости (требует sudo)"
	@echo "  make build         - Скомпилировать проект"
	@echo "  make run           - Запустить приложение (требует sudo)"
	@echo "  make clean         - Удалить бинарные файлы"

install-deps:
	@echo "Установка libpcap-dev..."
	sudo apt-get update
	sudo apt-get install -y libpcap-dev
	@echo "Установка Go зависимостей..."
	go mod download

build:
	@echo "Компиляция проекта..."
	go build -o arp-tool
	@echo "Готово! Запустите: sudo ./arp-tool"

run: build
	@echo "Запуск приложения (требует sudo)..."
	sudo ./arp-tool

clean:
	@echo "Удаление бинарных файлов..."
	rm -f arp-tool
	@echo "Готово!"
