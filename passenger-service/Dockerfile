FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o passenger-service ./cmd/passenger-service

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/passenger-service .

ENV PASSENGER_SERVICE_PORT=8082

EXPOSE 8082

CMD ["./passenger-service"]