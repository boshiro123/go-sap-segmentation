#!/bin/bash
set -e

# Переходим в корневую директорию проекта
cd "$(dirname "$0")/.."

# Проверяем наличие Go
if ! command -v go &> /dev/null; then
    echo "Go не установлен. Установите Go для сборки проекта."
    exit 1
fi

# Загружаем зависимости
echo "Загрузка зависимостей..."
go mod download

# Собираем проект
echo "Сборка проекта..."
go build -o bin/sap_segmentationd ./cmd/sap_segmentationd

# Проверяем успешность сборки
if [ -f bin/sap_segmentationd ]; then
    echo "Проект успешно собран! Исполняемый файл: bin/sap_segmentationd"
else
    echo "Ошибка при сборке проекта."
    exit 1
fi 