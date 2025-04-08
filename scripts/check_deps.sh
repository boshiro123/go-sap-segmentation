#!/bin/bash
set -e

echo "Проверка необходимых зависимостей..."

# Проверка наличия Docker
if ! command -v docker &> /dev/null; then
    echo "Docker не установлен. Установите Docker для запуска проекта."
    exit 1
fi

# Проверка наличия Docker Compose
if ! command -v docker-compose &> /dev/null; then
    echo "Docker Compose не установлен. Установите Docker Compose для запуска проекта."
    exit 1
fi

# Проверка наличия Go (требуется для разработки)
if ! command -v go &> /dev/null; then
    echo "Go не установлен. Для разработки рекомендуется установить Go."
    echo "Но вы можете запустить проект в Docker без локальной установки Go."
else
    go_version=$(go version | awk '{print $3}')
    echo "Используется Go версии: $go_version"
    
    # Проверка версии Go
    if [[ "$go_version" < "go1.18" ]]; then
        echo "Предупреждение: рекомендуется использовать Go версии 1.18 или выше."
    fi
fi

# Проверка наличия необходимых переменных окружения
env_file="./compose/.env"
if [ ! -f "$env_file" ]; then
    echo "Файл .env не найден в директории compose. Создайте его на основе примера."
    exit 1
fi

echo "✅ Все необходимые зависимости установлены." 