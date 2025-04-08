FROM golang:1.24-alpine AS builder

WORKDIR /build

# Копируем файлы go.mod и go.sum
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем исходный код
COPY . .

# Генерируем Swagger документацию
RUN go install github.com/swaggo/swag/cmd/swag@latest && \
    /go/bin/swag init -g cmd/sap_segmentationd/main.go -o docs/generated --ot json --parseInternal

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o sap_segmentationd ./cmd/sap_segmentationd

# Создаем минимальный образ
FROM alpine:3.19

WORKDIR /app

# Устанавливаем необходимые пакеты
RUN apk --no-cache add ca-certificates tzdata

# Копируем исполняемый файл из предыдущего этапа
COPY --from=builder /build/sap_segmentationd .

# Копируем Swagger документацию
COPY --from=builder /build/docs/generated /app/docs/generated

# Создаем директорию для логов
RUN mkdir -p /app/log

# Устанавливаем переменные окружения по умолчанию
ENV ENV=prod

# Запускаем приложение
CMD ["./sap_segmentationd"] 