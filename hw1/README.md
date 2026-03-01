# hw1

1. перейти в папку проекта
   cd <dir_name>

2. собрать проект mvn
   mvn clean compile

3. запустить сервер:
   mvn exec:java -Dexec.mainClass="ru.hse.network.Server" -Dexec.args="8080"

4. запустить клиент во втором терминале
   mvn exec:java -Dexec.mainClass="ru.hse.network.Client" -Dexec.args="127.0.0.1 8080 8 5000 25"

# сценарии

1. маленькие размеры
   java -cp target/classes ru.hse.network.Client 127.0.0.1 8080 8 5000 25 > results_scenario1.csv
2. большие размеры
   java -cp target/classes ru.hse.network.Client 127.0.0.1 8080 1024 5000 10 > results_scenario2.csv
3. свои параметры
   java -cp target/classes ru.hse.network.Client 127.0.0.1 8080 64 3000 20 > results_scenario3.csv

# запустить на виртуалке

1. на сервере:
   java -Xms100M -Xmx200M -cp target/network-throughput-1.0.jar ru.hse.network.Server 8080

2. на клиенте (замените IP на реальный IP сервера):
   java -cp target/network-throughput-1.0.jar ru.hse.network.Client <SERVER_IP> 8080 8 5000 25 > results.csv

# WireShark:

1. Запустить WireShark
2. Выбрать сетевой интерфейс
3. Применить фильтр: tcp.port == 8080
4. Запустить сервер и клиент
5. Сделать скриншоты пакетов

# Мониторинг системы:

- htop - для просмотра загрузки CPU и памяти
- nethogs - для просмотра сетевого трафика по процессам
