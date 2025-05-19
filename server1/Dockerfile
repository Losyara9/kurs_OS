FROM golang:1.22-alpine

WORKDIR /app

# Устанавливаем зависимости для xdotool и dmesg
RUN apk add --no-cache \
    xdotool \
    util-linux \
    linux-headers \
    musl-dev

COPY go.mod .
RUN go mod download

COPY . .
RUN go build -o server .

CMD ["./server"]
