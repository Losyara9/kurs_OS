package main

import (
	"fmt"
	"log"
	"net"
	"os/exec"
	"strings"
	"time"
)

func getCursorPosition() string {
	cmd := exec.Command("xdotool", "getmouselocation")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Sprintf("Cursor Error: %v", err)
	}
	return string(output)
}

func getLastError() string {
	cmd := exec.Command("dmesg", "-T", "-l", "err")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Sprintf("Error getting logs: %v", err)
	}
	lines := strings.Split(string(output), "\n")
	if len(lines) > 0 {
		return lines[len(lines)-1]
	}
	return "No errors found"
}

func handleConnection(conn net.Conn, logger chan string) {
	defer conn.Close()
	clientAddr := conn.RemoteAddr().String()
	logger <- fmt.Sprintf("Client connected: %s", clientAddr)

	lastSent := ""
	iteration := 0

	for {
		cursor := getCursorPosition()
		lastError := getLastError()
		errorMessage := ""

		iteration++
		if iteration%10 == 0 {
			errorMessage = fmt.Sprintf("FICTIVE ERROR: Simulated failure at %s", time.Now().Format("15:04:05"))
			logger <- errorMessage
		}

		// Формируем данные БЕЗ времени для сравнения
		dataWithoutTime := fmt.Sprintf("Cursor: %s\nLast Error: %s", cursor, lastError)
		if errorMessage != "" {
			dataWithoutTime += fmt.Sprintf("\nSynthetic Error: %s", errorMessage)
		}

		// Проверка на изменения по содержимому, не включая время
		if dataWithoutTime != lastSent {
			timestamp := time.Now().Format("2006-01-02 15:04:05")
			messageToSend := fmt.Sprintf("Time: %s\n%s", timestamp, dataWithoutTime)

			_, err := conn.Write([]byte(messageToSend + "\n"))
			if err != nil {
				logger <- fmt.Sprintf("Client %s disconnected", clientAddr)
				return
			}
			lastSent = dataWithoutTime
			logger <- fmt.Sprintf("Sent data to %s", clientAddr)
		}

		time.Sleep(1 * time.Second)
	}
}

func main() {

	time.Sleep(2 * time.Second)
	// Проверка уникальности сервера
	if _, err := net.Listen("tcp", ":18081"); err != nil {
		log.Fatal("Server1 is already running")
	}

	logger := make(chan string)
	go func() {
		for msg := range logger {
			log.Printf("SERVER1: %s", msg)
			// Отправка в сервер логирования
			if conn, err := net.Dial("tcp", "logger:8083"); err == nil {
				fmt.Fprintf(conn, "[server1] %s\n", msg)
				conn.Close()
			}
		}
	}()

	ln, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	logger <- "Server started on :8081"

	for {
		conn, err := ln.Accept()
		if err != nil {
			logger <- fmt.Sprintf("Accept error: %v", err)
			continue
		}
		go handleConnection(conn, logger)
	}
}
