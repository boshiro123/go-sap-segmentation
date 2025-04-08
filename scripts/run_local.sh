#!/bin/bash
set -e

# Переходим в корневую директорию проекта
cd "$(dirname "$0")/.."

# Проверяем наличие Go
if ! command -v go &> /dev/null; then
    echo "Go не установлен. Установите Go для локального запуска проекта."
    exit 1
fi

# Создаем директорию для логов, если она не существует
mkdir -p log

# Проверяем, запущен ли PostgreSQL
source compose/.env

if ! pg_isready -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" > /dev/null 2>&1; then
    echo "PostgreSQL не запущен или недоступен. Запустите базу данных перед запуском проекта."
    echo "Вы можете запустить только PostgreSQL командой:"
    echo "docker-compose -f compose/docker-compose.yml up -d postgres"
    exit 1
fi

# Собираем проект
echo "Сборка проекта..."
go build -o bin/sap_segmentationd ./cmd/sap_segmentationd

# Запускаем проект
echo "Запуск проекта в локальном режиме..."
ENV=local \
DB_HOST="$DB_HOST" \
DB_PORT="$DB_PORT" \
DB_NAME="$DB_NAME" \
DB_USER="$DB_USER" \
DB_PASSWORD="$DB_PASSWORD" \
CONN_URI="$CONN_URI" \
CONN_AUTH_LOGIN_PWD="$CONN_AUTH_LOGIN_PWD" \
CONN_USER_AGENT="$CONN_USER_AGENT" \
CONN_TIMEOUT="$CONN_TIMEOUT" \
CONN_INTERVAL="$CONN_INTERVAL" \
IMPORT_BATCH_SIZE="$IMPORT_BATCH_SIZE" \
LOG_CLEANUP_MAX_AGE="$LOG_CLEANUP_MAX_AGE" \
./bin/sap_segmentationd 