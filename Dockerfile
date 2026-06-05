# ================== Builder Stage ==================
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Копируем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Генерируем Swagger
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init -g cmd/app/main.go --parseDependency --parseInternal --parseDepth 1

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -v -o main ./cmd/app

# ================== Final Stage ==================
FROM golang:1.26-alpine

WORKDIR /app

# Копируем бинарник
COPY --from=builder /app/main .

# Копируем исходники (нужно для swag и тестов)
COPY --from=builder /app/go.mod /app/go.sum ./
COPY --from=builder /app/internal ./internal
COPY --from=builder /app/cmd ./cmd
COPY --from=builder /app/docs ./docs     

# Устанавливаем swag (чтобы можно было перегенерировать внутри контейнера)
RUN go install github.com/swaggo/swag/cmd/swag@latest && \
    go mod download

EXPOSE 8080

CMD ["./main"]