#!/bin/bash
set -e

# Переходим в корневую директорию проекта
cd "$(dirname "$0")/.."

# Проверяем, запущен ли контейнер с PostgreSQL
if ! docker ps | grep -q sap_segmentation_postgres; then
  echo "PostgreSQL не запущен. Запустите проект с помощью scripts/run.sh"
  exit 1
fi

# Получаем переменные окружения из файла .env
source compose/.env

echo "Применение миграций..."

# Запускаем миграцию с помощью psql в контейнере
docker exec -i sap_segmentation_postgres psql -U "$DB_USER" -d "$DB_NAME" < setup/install.sql

echo "Миграции успешно применены!" 