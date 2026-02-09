FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY authentication/go.mod authentication/go.sum ./authentication/
COPY proto ./proto

WORKDIR /app/authentication

RUN go mod download

COPY authentication .

RUN CGO_ENABLED=0 GOOS=linux go build -o auth ./cmd
FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/authentication .

CMD ["/bin/sh", "-c", "./auth > auth_output.log"]
