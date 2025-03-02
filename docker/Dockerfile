FROM golang:1.24.0-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy

COPY . .

RUN go build -o /app/bin/app ./cmd/main.go

FROM alpine:latest

WORKDIR /service

COPY --from=builder /app/bin/app .
COPY --from=builder /app/config.env /app/config.env

CMD ["./app"]
