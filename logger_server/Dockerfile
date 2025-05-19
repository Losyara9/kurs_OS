FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o logger . && chmod +x logger

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/logger .
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
RUN mkdir -p /app/logs
CMD ["/app/logger"]
