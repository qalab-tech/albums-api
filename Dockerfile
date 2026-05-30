FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

# Копируем весь код
COPY . .

# Собираем приложение (Swagger уже должен быть сгенерирован локально)
RUN CGO_ENABLED=0 GOOS=linux go build -v -o main ./cmd/app

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["./main"]