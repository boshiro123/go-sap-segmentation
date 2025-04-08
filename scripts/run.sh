#!/bin/bash
set -e

# Переходим в корневую директорию проекта
cd "$(dirname "$0")/.."

# Проверяем зависимости
echo "Проверка зависимостей..."
bash ./scripts/check_deps.sh

# Запускаем проект с помощью Docker Compose
echo "Запуск проекта..."
docker-compose -f compose/docker-compose.yml up -d

echo "Проект успешно запущен!"
echo "PostgreSQL доступен на порту $(grep DB_PORT compose/.env | cut -d '=' -f2)"
echo "Логи можно посмотреть командой: docker-compose -f compose/docker-compose.yml logs -f"
echo "Для остановки проекта используйте: bash scripts/stop.sh" 