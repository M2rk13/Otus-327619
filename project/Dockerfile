FROM golang:1.24.2-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -x -v -o currency-converter ./cmd

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/currency-converter .

EXPOSE 8080

CMD ["./currency-converter"]
