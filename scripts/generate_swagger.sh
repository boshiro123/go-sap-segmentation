#!/bin/bash
set -e

# Переходим в корневую директорию проекта
cd "$(dirname "$0")/.."

# Определяем версию swag для скачивания
SWAG_VERSION="v1.16.4"

echo "Генерация Swagger документации..."

# Используем go run для запуска swag без необходимости устанавливать его в систему
go run github.com/swaggo/swag/cmd/swag@$SWAG_VERSION init \
    -g cmd/sap_segmentationd/main.go \
    -o docs/generated \
    --ot json \
    --pd \
    --parseInternal

echo "Swagger документация успешно сгенерирована в директории docs/generated"
echo "Доступна по адресу: http://localhost:${APP_PORT:-8080}/swagger/index.html" 