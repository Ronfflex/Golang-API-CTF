package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env ? Source .env first ?")
	}
}

func CheckPort(host string, port int, timeout time.Duration) bool {
	address := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}

func FindOpenPorts(host string, start, end int, timeout time.Duration) []int {
	var wg sync.WaitGroup
	openPorts := []int{}
	portChan := make(chan int)

	for i := start; i <= end; i++ {
		wg.Add(1)
		go func(port int) {
			defer wg.Done()
			if CheckPort(host, port, timeout) {
				portChan <- port
			}
		}(i)
	}

	go func() {
		wg.Wait()
		close(portChan)
	}()

	for port := range portChan {
		openPorts = append(openPorts, port)
	}

	return openPorts
}

func main() {
	ip := os.Getenv("IP")
	startPortStr := os.Getenv("START_PORT")
	endPortStr := os.Getenv("END_PORT")
	timeoutMillisStr := os.Getenv("TIMEOUT")
	timeoutMillis, err := strconv.Atoi(timeoutMillisStr)
	if err != nil {
		log.Fatalf("Error converting timeout to integer: %v", err)
	}
	timeout := time.Duration(timeoutMillis) * time.Millisecond

	startPort, err := strconv.Atoi(startPortStr)
	if err != nil {
		log.Fatalf("Error converting start port to integer: %v", err)
	}

	endPort, err := strconv.Atoi(endPortStr)
	if err != nil {
		log.Fatalf("Error converting end port to integer: %v", err)
	}

	openPorts := FindOpenPorts(ip, startPort, endPort, timeout)
	if len(openPorts) > 0 {
		fmt.Println("OPEN PORTS FOUND : ", openPorts)
	} else {
		fmt.Println("NO OPEN PORTS FOUND !???")
	}
}
