#!/bin/bash
set -e

# Переходим в корневую директорию проекта
cd "$(dirname "$0")/.."

# Проверяем наличие Docker
if ! command -v docker &> /dev/null; then
    echo "Docker не установлен. Установите Docker для сборки проекта."
    exit 1
fi

# Проверяем наличие docker-compose
if ! command -v docker-compose &> /dev/null; then
    echo "docker-compose не установлен. Установите docker-compose для сборки проекта."
    exit 1
fi

# Генерируем Swagger документацию локально (если установлен swag)
if command -v swag &> /dev/null; then
    echo "Генерация Swagger-документации..."
    mkdir -p docs/generated
    swag init -g cmd/sap_segmentationd/main.go -o docs/generated --ot json --parseInternal || true
fi

# Собираем Docker-образ
echo "Сборка Docker-образа..."
docker-compose -f compose/docker-compose.yml build

# Проверяем успешность сборки
if [ $? -eq 0 ]; then
    echo "Docker-образ успешно собран!"
    echo "Теперь вы можете запустить проект командой: bash scripts/run.sh"
    echo "Для запуска с тестовыми данными, в файле compose/.env установите USE_TEST_DATA=true"
else
    echo "Ошибка при сборке Docker-образа."
    exit 1
fi 