#!/bin/bash
set -e

# Переходим в корневую директорию проекта
cd "$(dirname "$0")/.."

# Проверяем зависимости
echo "Проверка зависимостей..."
bash ./scripts/check_deps.sh

# Проверяем аргументы запуска
USE_TEST=true
if [ "$1" == "--test" ] || [ "$1" == "-t" ]; then
    USE_TEST=true
    echo "Включен режим использования тестовых данных"
fi

# Останавливаем текущие контейнеры, если они запущены
echo "Останавливаем текущие контейнеры, если они запущены..."
docker-compose -f compose/docker-compose.yml down

# Устанавливаем USE_TEST_DATA
if [ "$USE_TEST" == "true" ]; then
    export USE_TEST_DATA=true
    echo "USE_TEST_DATA=true установлен для текущего запуска"
fi

# Запускаем проект с помощью Docker Compose
echo "Запуск проекта..."
docker-compose -f compose/docker-compose.yml up -d

echo "Проект успешно запущен!"
echo "PostgreSQL доступен на порту $(grep DB_PORT compose/.env | cut -d '=' -f2)"
echo "Логи можно посмотреть командой: docker-compose -f compose/docker-compose.yml logs -f"
echo "Для остановки проекта используйте: bash scripts/stop.sh"

if [ "$USE_TEST" == "true" ]; then
    echo -e "\nВнимание: Используются тестовые данные вместо SAP API"
fi 