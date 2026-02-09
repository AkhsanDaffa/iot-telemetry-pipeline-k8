# Stage 1: Builder
FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod ./

COPY . .

RUN go build -o iot-app main.go

# Stage 2: Production 
FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/iot-app .

CMD ["./iot-app"]