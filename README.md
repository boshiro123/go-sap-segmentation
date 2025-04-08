# SAP Segmentation

Модуль для импорта данных сегментации из SAP API в PostgreSQL базу данных.

## Описание

Данный модуль выполняет следующие функции:

- Импортирует данные из внешнего SAP API в таблицу PostgreSQL
- Обновляет существующие данные при повторном импорте
- Логирует процесс импорта в консоль и файл
- Автоматически удаляет устаревшие логи
- Предоставляет REST API для доступа к данным и управления импортом
- Включает документацию API в формате Swagger
- Поддерживает работу с тестовыми данными при недоступности SAP API

## Требования

- Docker и Docker Compose
- Go 1.18+ (для локальной разработки)

## Структура проекта

```
.
├── build/                 # Файлы для сборки Docker образа
├── cmd/                   # Исполняемые файлы
│   └── sap_segmentationd/ # Основное приложение
├── compose/               # Docker Compose файлы и конфигурация
├── docs/                  # Документация API (Swagger)
│   └── generated/         # Автоматически сгенерированная документация
├── internal/              # Внутренние пакеты приложения
│   ├── api/               # API сервер
│   ├── logutil/           # Утилиты для работы с логами
│   ├── sap/               # Клиент для SAP API
│   └── storage/           # Работа с базой данных
├── model/                 # Модели данных
├── pkg/                   # Разделяемые пакеты
├── scripts/               # Скрипты для управления проектом
├── setup/                 # SQL миграции
└── log/                   # Директория для логов
```

## Быстрый старт

### Запуск с помощью Docker

1. Клонируйте репозиторий:

   ```bash
   git clone <repository-url>
   cd <repository-dir>
   ```

2. Запустите проект:

   ```bash
   ./scripts/run.sh
   ```

3. Откройте Swagger UI:

   ```
   http://localhost:8080/swagger/index.html
   ```

   Также доступен редирект с корневого пути:

   ```
   http://localhost:8080/
   ```

4. Для остановки проекта:
   ```bash
   ./scripts/stop.sh
   ```

### Локальная разработка

1. Установите зависимости:

   ```bash
   go mod download
   ```

2. Сгенерируйте Swagger документацию:

   ```bash
   ./scripts/generate_swagger.sh
   ```

3. Соберите проект:

   ```bash
   ./scripts/build.sh
   ```

4. Запустите PostgreSQL:

   ```bash
   docker-compose -f compose/docker-compose.yml up -d postgres
   ```

5. Запустите проект локально:
   ```bash
   ./scripts/run_local.sh
   ```

## API Endpoints

Проект предоставляет следующие REST API эндпоинты:

| Метод | Путь                     | Описание                              |
| ----- | ------------------------ | ------------------------------------- |
| GET   | /api/health              | Проверка работоспособности сервера    |
| GET   | /api/segmentation        | Получение всех сегментов              |
| GET   | /api/segmentation/:id    | Получение сегмента по SAP ID          |
| POST  | /api/segmentation/import | Запуск импорта сегментации из SAP API |
| GET   | /swagger/\*              | Документация API (Swagger UI)         |
| GET   | /                        | Редирект на Swagger UI                |

## Конфигурация

Конфигурация проекта осуществляется через переменные окружения. Значения по умолчанию:

| Переменная          | Значение по умолчанию                                        | Описание                            |
| ------------------- | ------------------------------------------------------------ | ----------------------------------- |
| DB_HOST             | 127.0.0.1                                                    | IP-адрес сервера БД                 |
| DB_PORT             | 5432                                                         | TCP-порт сервера БД                 |
| DB_NAME             | mesh_group                                                   | Название БД                         |
| DB_USER             | postgres                                                     | Имя пользователя БД                 |
| DB_PASSWORD         | postgres                                                     | Пароль пользователя БД              |
| CONN_URI            | http://bsm.api.iql.ru/ords/bsm/segmentation/get_segmentation | URL для подключения к внешнему API  |
| CONN_AUTH_LOGIN_PWD | 4Dfddf5:jKlljHGH                                             | Логин и пароль для аутентификации   |
| CONN_USER_AGENT     | spacecount-test                                              | User-Agent для подключения к SAP    |
| CONN_TIMEOUT        | 5s                                                           | Таймаут подключения к внешнему API  |
| CONN_INTERVAL       | 1500ms                                                       | Задержка между запросами            |
| IMPORT_BATCH_SIZE   | 50                                                           | Размер пачки данных при запросе     |
| LOG_CLEANUP_MAX_AGE | 7                                                            | Время хранения логов в днях         |
| APP_PORT            | 8080                                                         | Порт для HTTP сервера               |
| RUN_IMPORT_ON_START | true                                                         | Запускать импорт при старте сервера |
