FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o driver-service ./cmd/driver-service

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/driver-service .

ENV DRIVER_SERVICE_PORT=8081

EXPOSE 8081

CMD ["./driver-service"]