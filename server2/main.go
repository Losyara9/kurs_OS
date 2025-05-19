package main

import (
	"fmt"
	"log"
	"net"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func getMemoryUsage() (float64, float64) {
	cmd := exec.Command("free")
	output, err := cmd.Output()
	if err != nil {
		return 0, 0
	}

	lines := strings.Split(string(output), "\n")
	if len(lines) < 2 {
		return 0, 0
	}

	// Физическая память
	fields := strings.Fields(lines[1])
	totalPhys, _ := strconv.ParseFloat(fields[1], 64)
	usedPhys, _ := strconv.ParseFloat(fields[2], 64)
	physPercent := (usedPhys / totalPhys) * 100

	// Виртуальная память
	if len(lines) > 2 {
		fields = strings.Fields(lines[2])
		totalVirt, _ := strconv.ParseFloat(fields[1], 64)
		usedVirt, _ := strconv.ParseFloat(fields[2], 64)
		virtPercent := (usedVirt / totalVirt) * 100
		return physPercent, virtPercent
	}

	return physPercent, 0
}

func handleConnection(conn net.Conn, logger chan string) {
	defer conn.Close()
	clientAddr := conn.RemoteAddr().String()
	logger <- fmt.Sprintf("Client connected: %s", clientAddr)

	lastPhys, lastVirt := -1.0, -1.0
	for {
		phys, virt := getMemoryUsage()

		// Проверка ИМЕННО на изменения физической или виртуальной памяти
		if phys != lastPhys || virt != lastVirt {
			timestamp := time.Now().Format("2006-01-02 15:04:05")
			msg := fmt.Sprintf("Time: %s\nPhysical: %.1f%%\nVirtual: %.1f%%", timestamp, phys, virt)

			_, err := conn.Write([]byte(msg + "\n"))
			if err != nil {
				logger <- fmt.Sprintf("Client %s disconnected", clientAddr)
				return
			}
			lastPhys, lastVirt = phys, virt
			logger <- fmt.Sprintf("Sent to %s: %s", clientAddr, msg)
		}

		time.Sleep(1 * time.Second)
	}
}

func main() {

	time.Sleep(2 * time.Second)
	// Проверка уникальности сервера
	if _, err := net.Listen("tcp", ":18082"); err != nil {
		log.Fatal("Server2 is already running")
	}

	logger := make(chan string)
	go func() {
		for msg := range logger {
			log.Printf("SERVER2: %s", msg)
			// Отправка в сервер логирования
			if conn, err := net.Dial("tcp", "logger:8083"); err == nil {
				fmt.Fprintf(conn, "[server2] %s\n", msg)
				conn.Close()
			}
		}
	}()

	ln, err := net.Listen("tcp", ":8082")
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	logger <- "Server started on :8082"

	for {
		conn, err := ln.Accept()
		if err != nil {
			logger <- fmt.Sprintf("Accept error: %v", err)
			continue
		}
		go handleConnection(conn, logger)
	}
}
