package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

func handleConnection(conn net.Conn, logFile *os.File) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		msg := fmt.Sprintf("[%s] %s\n", time.Now().Format("2006-01-02 15:04:05"), scanner.Text())
		log.Print(msg)
		if _, err := logFile.WriteString(msg); err != nil {
			log.Printf("Write error: %v", err)
		}
	}
}

func main() {
	// Создаем директорию для логов если её нет
	if err := os.MkdirAll("/app/logs", 0755); err != nil {
		log.Fatalf("Failed to create logs directory: %v", err)
	}

	// Создаем файл логов
	logFile, err := os.OpenFile("/app/logs/server.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to create log file: %v", err)
	}
	defer logFile.Close()

	// Проверка записи
	if _, err := logFile.WriteString("Logger initialized at " + time.Now().Format(time.RFC3339) + "\n"); err != nil {
		log.Fatalf("Can't write to log file: %v", err)
	}

	ln, err := net.Listen("tcp", ":8083")
	if err != nil {
		logFile.WriteString("Failed to start server: " + err.Error() + "\n")
		log.Fatal(err)
	}
	defer ln.Close()

	log.Printf("Logger server started on :8083")
	logFile.WriteString("Server started\n")

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("Accept error: %v", err)
			continue
		}
		go handleConnection(conn, logFile)
	}
}
