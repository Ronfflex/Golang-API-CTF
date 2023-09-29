package main

import (
	"bytes"
	"encoding/json"
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

type PostBody struct {
	User   string `json:"User,omitempty"`
	Secret string `json:"Secret,omitempty"`
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

func FetchDetails(ip string, port int, path string) ResponseDetails {
	url := fmt.Sprintf("http://%s:%d%s", ip, port, path)
	resp, err := http.Get(url)
	if err != nil {
		return ResponseDetails{Error: err.Error()}
	}

	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return ResponseDetails{Error: err.Error()}
	}

	return ResponseDetails{
		Status:  resp.Status,
		Headers: resp.Header,
		Body:    string(bodyBytes),
	}
}

func postBodyToCheckReponse(ip string, port int, path string, body PostBody) ([]byte, error) {
	url := fmt.Sprintf("http://%s:%d%s", ip, port, path)

	jsonBody, err := json.Marshal(body)
	if err != nil {
		fmt.Println("Error marshaling body:", err)
		return nil, err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Println("Error making POST request:", err)
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading POST response body:", err)
		return nil, err
	}

	//fmt.Println("[EMPTY]", path, "Response:", string(respBody))
	return respBody, nil
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

	paths := []string{"/ping", "/signup", "/check", "/getUserSecret", "/getUserLevel", "/getUserPoints", "/iNeedAHint", "/enterChallenge", "/submitSolution"}

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

	openPorts := FindOpenPorts(ip, startPort, endPort, timeout)
	if len(openPorts) == 0 {
		fmt.Println("Couldn't find any open ports in the specified range.")
		return
	} else {
		fmt.Println("Found open ports:", openPorts)
	}

	for _, port := range openPorts {
		for _, path := range paths {
			fmt.Println("Fetching details for port :", port, ", path", path)
			details := FetchDetails(ip, port, path)

			if details.Error != "" {
				fmt.Println("Port,", port, "Path", path, ":\n Error :", details.Error)
				break
			}

			fmt.Println("------------------GETTING INFOS---------------------")
			fmt.Println("Port", port, ": Status :", details.Status)
			fmt.Println("Port", port, ": Headers :", details.Headers)
			fmt.Println("Port", port, ": Body :", details.Body)
			//fmt.Println("Making POST request with empty body...")
			//postBodyToCheckReponse(ip, port, path, PostBody{}) // Empty body to check response

			user := "testUser"
			var secret string

			for _, path := range paths {
				postBody := PostBody{User: user}
				if path == "/getUserLevel" || path == "/getUserPoints" {
					postBody.Secret = secret
					fmt.Println(postBody.Secret)
				}

				fmt.Println("----------------POST ON", path, "-----------------")
				respBody, err := postBodyToCheckReponse(ip, port, path, postBody)
				if err != nil {
					fmt.Println("Error:", err)
					break
				}

				if path == "/getUserSecret" {
					secret = string(respBody)
				}

				fmt.Println(path, "Response:", string(respBody))
			}
		}
	}
}
