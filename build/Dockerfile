FROM golang:1.24-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go install github.com/swaggo/swag/cmd/swag@latest && \
    /go/bin/swag init -g cmd/sap_segmentationd/main.go -o docs/generated --ot json --parseInternal

RUN CGO_ENABLED=0 GOOS=linux go build -o sap_segmentationd ./cmd/sap_segmentationd

FROM alpine:3.19

WORKDIR /app

RUN apk --no-cache add ca-certificates tzdata

COPY --from=builder /build/sap_segmentationd .

COPY --from=builder /build/docs/generated /app/docs/generated

RUN mkdir -p /app/log

ENV ENV=prod

CMD ["./sap_segmentationd"] 