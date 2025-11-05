FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init -g cmd/app/main.go --output docs

RUN go build -o petelka-api cmd/appy/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/petelka-api .
COPY --from=builder /app/docs ./docs
COPY .env .

EXPOSE 8080

CMD ["./petelka-api"]