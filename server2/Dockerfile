FROM golang:1.22-alpine

WORKDIR /app

# Устанавливаем free для мониторинга памяти
RUN apk add --no-cache procps

COPY go.mod .
RUN go mod download

COPY . .
RUN go build -o server .

CMD ["./server"]
