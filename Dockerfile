FROM golang:1.21-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init -g cmd/main.go

RUN go build -o ecommerce cmd/main.go

EXPOSE 8080

CMD ["./ecommerce"]