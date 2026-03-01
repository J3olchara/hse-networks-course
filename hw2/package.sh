#!/bin/bash

# Скрипт для упаковки проекта в ZIP архив для сдачи HW2

echo "=== Упаковка проекта HW2 ==="
echo ""

# Запросить фамилию
read -p "Введите вашу фамилию (для названия архива): " SURNAME

if [ -z "$SURNAME" ]; then
    echo "Ошибка: фамилия не может быть пустой"
    exit 1
fi

echo ""
echo "Создание архива: ${SURNAME}.zip"

# Создать временную директорию
TEMP_DIR="/tmp/hw2_package_$$"
mkdir -p "$TEMP_DIR/$SURNAME"

# Копировать файлы проекта
echo "Копирование файлов проекта..."
cp -r ../hw2/* "$TEMP_DIR/$SURNAME/"

# Удалить ненужные файлы
echo "Очистка..."
rm -f "$TEMP_DIR/$SURNAME/.env"
rm -f "$TEMP_DIR/$SURNAME/arp-tool"
rm -f "$TEMP_DIR/$SURNAME/package.sh"

# Создать ZIP архив
cd "$TEMP_DIR"
zip -r "${SURNAME}.zip" "$SURNAME"

# Переместить архив в текущую директорию
mv "${SURNAME}.zip" -
cd -

# Удалить временную директорию
rm -rf "$TEMP_DIR"

echo ""
echo "Готово! Создан архив: ${SURNAME}.zip"
echo ""
echo "Структура архива:"
unzip -l "${SURNAME}.zip" | head -20

echo ""
echo "ВАЖНО: Не забудьте добавить в архив файл отчета!"
echo "Команда: zip ${SURNAME}.zip отчет.pdf"
echo "или:     zip ${SURNAME}.zip отчет.docx"
