#!/bin/bash
set -e

# Переходим в корневую директорию проекта
cd "$(dirname "$0")/.."

# Останавливаем проект с помощью Docker Compose
echo "Остановка проекта..."
docker-compose -f compose/docker-compose.yml stop

echo "Проект остановлен." 