FROM golang:1.24-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY . .

# Убираем автоматическую генерацию Swagger
# RUN go install github.com/swaggo/swag/cmd/swag@latest && \
#     /go/bin/swag init -g cmd/sap_segmentationd/main.go -o docs/generated --ot json --parseInternal

RUN CGO_ENABLED=0 GOOS=linux go build -o sap_segmentationd ./cmd/sap_segmentationd

FROM alpine:3.19

WORKDIR /app

RUN apk --no-cache add ca-certificates tzdata

COPY --from=builder /build/sap_segmentationd .

# Создаем необходимые директории
RUN mkdir -p /app/docs/generated /app/log

ENV ENV=prod
ENV USE_TEST_DATA=true

CMD ["./sap_segmentationd"] 