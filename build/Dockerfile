FROM golang:1.24-alpine AS builder

WORKDIR /build

# Копируем файлы go.mod и go.sum
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o sap_segmentationd ./cmd/sap_segmentationd

# Создаем минимальный образ
FROM alpine:3.19

WORKDIR /app

# Устанавливаем необходимые пакеты
RUN apk --no-cache add ca-certificates tzdata

# Копируем исполняемый файл из предыдущего этапа
COPY --from=builder /build/sap_segmentationd .

# Создаем директорию для логов
RUN mkdir -p /app/log

# Устанавливаем переменные окружения по умолчанию
ENV ENV=prod

# Запускаем приложение
CMD ["./sap_segmentationd"] 