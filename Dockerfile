FROM golang:1.24-alpine AS builder

WORKDIR /app

# Копируем go.mod и go.sum и загружаем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем остальной исходный код
COPY . .

COPY .env ./.env

# Устанавливаем swag
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Генерируем документацию Swagger
RUN swag init -g cmd/app/main.go

# Собираем приложение
RUN go build -o ecommerce cmd/app/main.go

# Этап запуска
FROM alpine:latest

# Устанавливаем ca-certificates для HTTPS запросов (важно для внешних API или баз данных)
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Копируем исполняемый файл из этапа сборки
COPY --from=builder /app/ecommerce .

# Копируем сгенерированные файлы Swagger UI (HTML, JS, CSS)
COPY --from=builder /app/docs ./docs

# Если .env нужен при запуске, он должен быть скопирован в финальный образ
COPY --from=builder /app/.env ./.env

EXPOSE 8080

CMD ["./ecommerce"]