package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
)

var (
	logFiles = make(map[string]*os.File)
	mu       sync.Mutex
)

func getLogFile(serverTag string) (*os.File, error) {
	mu.Lock()
	defer mu.Unlock()

	if file, exists := logFiles[serverTag]; exists {
		return file, nil
	}

	err := os.MkdirAll("logs", os.ModePerm)
	if err != nil {
		return nil, err
	}

	// Открываем файл логов для server1 или server2
	filename := fmt.Sprintf("logs/%s.log", serverTag)
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}

	logFiles[serverTag] = file
	return file, nil
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 1024)

	for {
		n, err := conn.Read(buf)
		if err != nil {
			return
		}

		msg := string(buf[:n])
		tag := extractServerTag(msg)
		if tag == "" {
			tag = "unknown"
		}

		file, err := getLogFile(tag)
		if err != nil {
			log.Println("error opening log file:", err)
			continue
		}

		file.WriteString(msg + "\n")
	}
}

func extractServerTag(msg string) string {
	// Ищем [server1] или [server2]
	if strings.Contains(msg, "[server1]") {
		return "server1"
	} else if strings.Contains(msg, "[server2]") {
		return "server2"
	}
	return ""
}

func main() {
	ln, err := net.Listen("tcp", ":8083")
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	log.Println("Log server started on port 9000")

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		go handleConnection(conn)
	}
}
