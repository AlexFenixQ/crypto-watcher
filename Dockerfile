FROM golang:1.20 AS builder

WORKDIR /app
COPY . .
# внутри stage-1 (final image)
COPY ./frontend /app/frontend

# Установка swag внутри контейнера
RUN go install github.com/swaggo/swag/cmd/swag@latest
ENV PATH="/go/bin:$PATH"

# Генерация Swagger-доков
RUN swag init

# Стабильные зависимости
RUN go mod tidy

ENV CGO_ENABLED=0 GOOS=linux
RUN go build -o app main.go

# FINAL STAGE
FROM alpine:latest
WORKDIR /root/
RUN apk --no-cache add ca-certificates

# Копируем app
COPY --from=builder /app/app .

# 👇 Копируем UI внутрь контейнера
COPY ./frontend /app/frontend

CMD ["./app"]