package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

type ResponseDetails struct {
	Status     string
	StatusCode int
	Headers    map[string][]string
	Body       string
	Error      string
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

func FetchDetails(ip string, port int, path string) *ResponseDetails {
	responseDetails := &ResponseDetails{}
	url := fmt.Sprintf("http://%s:%d%s", ip, port, path)

	resp, err := http.Get(url)
	if err != nil {
		responseDetails.Error = err.Error()
		return responseDetails
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		responseDetails.Error = err.Error()
		return responseDetails
	}

	responseDetails.Status = resp.Status
	responseDetails.Headers = resp.Header
	responseDetails.Body = string(body)
	return responseDetails
}


func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("No .env found")
	}

	ip := os.Getenv("IP")
	startPortStr := os.Getenv("START_PORT")
	endPortStr := os.Getenv("END_PORT")
	timeoutMillisStr := os.Getenv("TIMEOUT")
	timeoutMillis, err := strconv.Atoi(timeoutMillisStr)

	if err != nil {
		log.Fatal("Error converting timeout to integer: ", err)
	}
	timeout := time.Duration(timeoutMillis) * time.Millisecond

	startPort, err := strconv.Atoi(startPortStr)
	if err != nil {
		log.Fatal("Error converting start port to integer: ", err)
	}

	endPort, err := strconv.Atoi(endPortStr)
	if err != nil {
		log.Fatal("Error converting end port to integer: ", err)
	}

	paths := []string{"/ping", "/signup", "/check"}

	openPorts := FindOpenPorts(ip, startPort, endPort, timeout)
	if len(openPorts) == 0 {
		fmt.Println("Couldn't find any open ports in the specified range.")
		return
	}

	for _, port := range openPorts {
		for _, path := range paths {
			fmt.Println("Fetching details for port :", port,", path", path)
			details := FetchDetails(ip, port, path)

			if details.Error != "" {
				fmt.Println("Port,", port, "Path", path, ":\n Error :", details.Error)
				continue
			}

			fmt.Println("Port", port,": Status :", details.Status)
			fmt.Println("Port", port,": Headers :", details.Headers)
			fmt.Println("Port", port,": Body :", details.Body)
			fmt.Println("---------------------------------------------")
		}
	}
}
