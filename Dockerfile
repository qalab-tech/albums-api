# ================== Builder Stage ==================
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Копируем go.mod и go.sum
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь исходный код
COPY . .

# Генерируем Swagger
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init -g cmd/app/main.go --parseDependency --parseInternal --parseDepth 1

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -v -o main ./cmd/app

# ================== Final Stage (с Go для тестов) ==================
FROM golang:1.26-alpine

WORKDIR /app

# Копируем скомпилированное приложение
COPY --from=builder /app/main .

# Копируем исходники для возможности запуска тестов
COPY --from=builder /app/go.mod /app/go.sum ./
COPY --from=builder /app/internal ./internal
COPY --from=builder /app/cmd ./cmd

# Скачиваем зависимости
RUN go mod download

EXPOSE 8080

CMD ["./main"]