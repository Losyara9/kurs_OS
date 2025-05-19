package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"sync"
)

var (
	currentConns []net.Conn
	connMutex    sync.Mutex
)

func connectToServer(server string, output chan<- string) {
	conn, err := net.Dial("tcp", server)
	if err != nil {
		output <- fmt.Sprintf("Connection error (%s): %v", server, err)
		return
	}

	// Сохраняем подключение
	connMutex.Lock()
	currentConns = append(currentConns, conn)
	connMutex.Unlock()

	output <- fmt.Sprintf("Connected to %s", server)
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		output <- fmt.Sprintf("[%s] %s", server, scanner.Text())
	}
}

func closeAllConnections() {
	connMutex.Lock()
	for _, conn := range currentConns {
		conn.Close() // Это вызовет завершение scanner.Scan()
	}
	currentConns = nil
	connMutex.Unlock()
}

func main() {
	servers := map[int]string{
		1: "server1:8081",
		2: "server2:8082",
	}

	output := make(chan string)
	go func() {
		for msg := range output {
			fmt.Println(msg)
		}
	}()

	for {
		fmt.Println("\nChoose server:")
		fmt.Println("1. Server 1 (Cursor & Errors)")
		fmt.Println("2. Server 2 (Memory Usage)")
		fmt.Println("3. Both servers")
		fmt.Println("4. Exit")
		fmt.Print("> ")

		var choice int
		_, err := fmt.Scan(&choice)
		if err != nil {
			fmt.Println("Invalid input")
			continue
		}

		// Закрываем предыдущие соединения
		closeAllConnections()

		switch choice {
		case 1:
			go connectToServer(servers[1], output)
		case 2:
			go connectToServer(servers[2], output)
		case 3:
			go connectToServer(servers[1], output)
			go connectToServer(servers[2], output)
		case 4:
			closeAllConnections()
			fmt.Println("Exiting.")
			os.Exit(0)
		default:
			fmt.Println("Invalid choice")
		}
	}
}

