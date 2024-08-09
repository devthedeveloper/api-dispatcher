package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/quic-go/quic-go/http3"
)

type RequestConfig struct {
	URL     string            `json:"url"`
	Method  string            `json:"method"`
	Headers map[string]string `json:"headers"`
	Body    map[string]string `json:"body,omitempty"`
}

type Config struct {
	Requests []RequestConfig `json:"requests"`
}

func loadConfigFromBody(body []byte) (Config, error) {
	var config Config
	err := json.Unmarshal(body, &config)
	return config, err
}

func sendRequest(config RequestConfig, wg *sync.WaitGroup, results chan<- string) {
	defer wg.Done()

	client := &http.Client{}
	var req *http.Request
	var err error

	if config.Method == "POST" && config.Body != nil {
		bodyData, _ := json.Marshal(config.Body)
		req, err = http.NewRequest(config.Method, config.URL, bytes.NewBuffer(bodyData))
	} else {
		req, err = http.NewRequest(config.Method, config.URL, nil)
	}

	if err != nil {
		results <- fmt.Sprintf("Error creating request for %s: %v", config.URL, err)
		return
	}

	for key, value := range config.Headers {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		results <- fmt.Sprintf("Error sending request to %s: %v", config.URL, err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		results <- fmt.Sprintf("Error reading response from %s: %v", config.URL, err)
		return
	}

	result := fmt.Sprintf("Response from %s: %s", config.URL, string(body))
	results <- result
}

func handleAPIRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is supported", http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	config, err := loadConfigFromBody(body)
	if err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	var wg sync.WaitGroup
	results := make(chan string)

	go func() {
		for result := range results {
			fmt.Fprintln(w, result)
		}
	}()

	for _, reqConfig := range config.Requests {
		wg.Add(1)
		go sendRequest(reqConfig, &wg, results)
	}

	wg.Wait()
	close(results)
}

func main() {
	useHTTP3 := flag.Bool("http3", false, "Start an HTTP/3 server")
	addr := flag.String("addr", ":8080", "HTTP/1.1 and HTTP/2 server address")
	http3Addr := flag.String("http3-addr", ":8443", "HTTP/3 server address")
	configFile := flag.String("config", "", "Path to the configuration file")
	flag.Parse()

	certPath := "cert.pem"
	keyPath := "key.pem"

	if *useHTTP3 {
		// Debugging output
		fmt.Println("Using cert.pem path:", certPath)
		fmt.Println("Using key.pem path:", keyPath)

		if _, err := os.Stat(certPath); os.IsNotExist(err) {
			log.Fatalf("cert.pem file does not exist at: %s", certPath)
		}

		if _, err := os.Stat(keyPath); os.IsNotExist(err) {
			log.Fatalf("key.pem file does not exist at: %s", keyPath)
		}

		// Start HTTP/1.1 and HTTP/2 server
		go func() {
			server := &http.Server{
				Addr:    *addr,
				Handler: http.HandlerFunc(handleAPIRequest),
			}
			log.Printf("Starting HTTP/1.1 and HTTP/2 server on %s", *addr)
			log.Fatal(server.ListenAndServeTLS(certPath, keyPath))
		}()

		// Start HTTP/3 server
		tlsCert, err := tls.LoadX509KeyPair(certPath, keyPath)
		if err != nil {
			log.Fatalf("Failed to load TLS certificates: %v", err)
		}

		quicServer := &http3.Server{
			Addr:      *http3Addr,
			TLSConfig: &tls.Config{Certificates: []tls.Certificate{tlsCert}},
		}

		log.Printf("Starting HTTP/3 server on %s", *http3Addr)
		log.Fatal(quicServer.ListenAndServeTLS(certPath, keyPath)) // Use correct paths here
	} else if *configFile != "" {
		// Run as a CLI tool
		configData, err := ioutil.ReadFile(*configFile)
		if err != nil {
			log.Fatalf("Error reading config file: %v", err)
		}

		config, err := loadConfigFromBody(configData)
		if err != nil {
			log.Fatalf("Error loading config: %v", err)
		}

		var wg sync.WaitGroup
		results := make(chan string)

		go func() {
			for result := range results {
				fmt.Println(result)
			}
		}()

		for _, reqConfig := range config.Requests {
			wg.Add(1)
			go sendRequest(reqConfig, &wg, results)
		}

		wg.Wait()
		close(results)
	} else {
		log.Fatal("Usage: api-dispatcher -config=<config-file> or api-dispatcher -http3")
	}
}
