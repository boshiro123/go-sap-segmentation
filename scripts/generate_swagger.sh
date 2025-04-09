#!/bin/bash
set -e

# Переходим в корневую директорию проекта
cd "$(dirname "$0")/.."

# Проверяем, установлен ли swag
if ! command -v swag &> /dev/null; then
    echo "Утилита swag не установлена. Устанавливаю..."
    go install github.com/swaggo/swag/cmd/swag@latest
fi

# Создаем директорию для документации, если она не существует
mkdir -p docs/generated

# Генерируем Swagger документацию
echo "Генерация Swagger-документации..."
swag init -g cmd/sap_segmentationd/main.go -o docs/generated --ot json --parseInternal

echo "Swagger-документация успешно сгенерирована в docs/generated/"
echo "Доступна по адресу: http://localhost:${APP_PORT:-8080}/swagger/index.html"

# Подсказка для запуска Docker
echo -e "\nТеперь вы можете собрать Docker-образ командой:"
echo "docker-compose -f compose/docker-compose.yml build"
echo "И запустить проект командой:"
echo "docker-compose -f compose/docker-compose.yml up -d" 