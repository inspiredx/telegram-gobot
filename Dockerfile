FROM golang:1.22.5-alpine3.19 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o bot .


FROM alpine:3.19
RUN apk update && apk upgrade libcrypto3 libssl3 && rm -rf /var/cache/apk/*
WORKDIR /app
COPY --from=builder /app/bot .
CMD ["./bot"]
