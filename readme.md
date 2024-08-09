# Multi-threaded API Request Dispatcher

This project is a multi-threaded API request dispatcher written in Go. It supports HTTP/1.1, HTTP/2, and HTTP/3 and can be used as a command-line tool or as a server that processes API requests concurrently. This tool is ideal for testing load, simulating traffic, or gathering data from multiple APIs simultaneously.

## Features

- **Multi-threaded:** Concurrently dispatches multiple API requests using Go's goroutines.
- **Supports HTTP/1.1, HTTP/2, and HTTP/3:** Can handle HTTP requests across different protocols.
- **Flexible Configuration:** Configurable via a JSON file specifying the API endpoints, methods, headers, and bodies.
- **Optional HTTP/3 Server:** The server can be started to listen for HTTP/3 requests, along with HTTP/1.1 and HTTP/2.

## Prerequisites

- **Go (Golang):** You need Go installed on your machine to compile and run the program. Download it from [golang.org](https://golang.org/dl/).
- **OpenSSL:** Required for generating self-signed TLS certificates if you don't have valid certificates.

## Installation

### 1. Clone the Repository

```bash
git clone https://github.com/yourusername/api-dispatcher.git
cd api-dispatcher
```

### 2. Install Dependencies

Ensure the required Go packages are installed:

```bash
go get github.com/quic-go/quic-go/http3
```

### 3. Generate TLS Certificates (Optional)

If you don't have TLS certificates, generate self-signed certificates using OpenSSL:

```bash
openssl req -new -newkey rsa:2048 -days 365 -nodes -x509 -keyout key.pem -out cert.pem
```

This will create `cert.pem` and `key.pem` in your project directory.

## Compilation

### Compile the Go Program

You can compile the Go program into a binary executable:

```bash
go build -o api-dispatcher dispatcher.go
```

This will create an executable named `api-dispatcher`.

## Running the Application

### 1. Running as a CLI Tool

You can use the tool directly from the command line to send API requests based on a configuration file.

```bash
./api-dispatcher -config=config.json
```

- **`-config=config.json`**: Specifies the path to the JSON configuration file.

### 2. Running as an HTTP/1.1, HTTP/2, and HTTP/3 Server

You can start the server to listen for HTTP requests. The server will handle HTTP/1.1, HTTP/2, and HTTP/3 connections.

```bash
./api-dispatcher -http3
```

- **`-http3`**: Starts the server on ports `8080` (HTTP/1.1 and HTTP/2) and `8443` (HTTP/3).
- **`-addr=:8080`**: (Optional) Specify a different address or port for the HTTP/1.1 and HTTP/2 server.
- **`-http3-addr=:8443`**: (Optional) Specify a different address or port for the HTTP/3 server.

### 3. Configuration File Format

The configuration file should be in JSON format and specify the API requests. Here’s an example `config.json`:

```json
{
  "requests": [
    {
      "url": "https://jsonplaceholder.typicode.com/posts/1",
      "method": "GET",
      "headers": {
        "Accept": "application/json"
      }
    },
    {
      "url": "https://jsonplaceholder.typicode.com/posts",
      "method": "POST",
      "headers": {
        "Content-Type": "application/json"
      },
      "body": {
        "title": "foo",
        "body": "bar",
        "userId": "1"
      }
    }
  ]
}
```

## Testing

### 1. Testing HTTP/1.1 and HTTP/2

You can use cURL or Postman to test the HTTP/1.1 and HTTP/2 endpoints.

#### Using cURL

```bash
curl -v https://localhost:8080/dispatch -d @config.json --insecure
```

### 2. Testing HTTP/3

To test the HTTP/3 server, use cURL with the `--http3` flag.

#### Using cURL

```bash
curl --http3 -v https://localhost:8443/dispatch -d @config.json --insecure
```

### 3. Testing with a Web Browser

You can test the HTTP/3 server using a modern web browser like Chrome or Firefox:

```
https://localhost:8443/dispatch
```

## Integration with Other Languages

### 1. Integration with Python

You can run the Go binary from a Python script using the `subprocess` module.

```python
import subprocess

def run_api_dispatcher(config_path):
    result = subprocess.run(
        ["./api-dispatcher", "-config", config_path],
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
        text=True
    )
    print(result.stdout)

run_api_dispatcher("config.json")
```

### 2. Integration with Node.js

You can also execute the Go binary from a Node.js application using the `child_process` module.

```javascript
const { spawn } = require('child_process');

function runApiDispatcher(configPath) {
    const process = spawn('./api-dispatcher', ['-config', configPath]);

    process.stdout.on('data', (data) => {
        console.log(`Output: ${data}`);
    });

    process.stderr.on('data', (data) => {
        console.error(`Error: ${data}`);
    });

    process.on('close', (code) => {
        console.log(`Process exited with code: ${code}`);
    });
}

runApiDispatcher('config.json');
```

## Troubleshooting

- **No Such File or Directory Error:** Ensure that `cert.pem` and `key.pem` exist in the directory where you are running the Go binary.
- **Permission Issues:** Ensure that your Go binary has read access to the TLS certificate and key files.
- **Self-Signed Certificate Warnings:** When using cURL or Postman, add the `--insecure` flag or disable SSL verification to bypass warnings about self-signed certificates.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---

This README provides a detailed overview of how to set up, compile, run, and integrate the Go-based API dispatcher with various tools and languages. Let me know if you need any more details or adjustments!