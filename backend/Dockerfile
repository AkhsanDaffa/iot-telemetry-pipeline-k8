# Stage 1: Builder
FROM golang:1.23-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o iot-app main.go

# Stage 2: Production 
FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /root/

COPY --from=builder /app/iot-app .

EXPOSE 8080

CMD ["./iot-app"]