FROM golang:1.22-alpine

WORKDIR /app

COPY go.mod .
RUN go mod download

COPY . .
RUN go build -o client .

CMD ["sh", "-c", "exec ./client"]
