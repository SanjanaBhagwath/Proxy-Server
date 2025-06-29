# Go-Based Secure Proxy Server

This project is a modular, high-performance HTTP proxy server implemented in **Go**, designed to demonstrate backend concepts like request interception, IP masking, rate limiting, and malicious URL blocking using the **Google Safe Browsing API**. The system is fully containerized using **Docker**, with separate containers for each core component to simulate a scalable, real-world network.

## Why Go?

Go (Golang) is well-suited for building networked services due to its:

- Lightweight concurrency model using goroutines
- Fast compilation and execution
- Rich standard library for HTTP, networking, and cryptography
- Simplicity and readability, making it ideal for scalable microservices

These features make Go a natural fit for implementing a performant and modular proxy server.

## Architecture Overview

The system is composed of four Docker containers, each representing a distinct component in the proxy chain:

| Container           | Role                                                                 |
|---------------------|----------------------------------------------------------------------|
| `proxy-server`      | Main HTTP proxy that forwards requests and masks client IPs         |
| `target-server`     | Simulated backend server that responds to forwarded requests         |
| `rate-limiter`      | Middleware proxy that limits requests per client IP                 |
| `malicious-blocker` | Middleware proxy that blocks requests to malicious URLs using the Google Safe Browsing API |

The project was built in increments with each increment implementing an additional feature.

## Running the Project

### 1. Running the Client

Build and run Client.go

```bash
go build Client.go
go run Client.go
```

### 2. Building & Running the Proxy and Target Server

Ensure Docker is installed and running. Then build and run the containers on separate terminals

```bash
docker build -t proxy-server .
docker build -t target-server .

docker run --name TargetServer --network proxy-network -p 8080:8080 targetserver
docker run --name ProxyServer --network proxy-network -p 9090:9090 proxyserver
```
To test the target server connection, on the client terminal, enter the URL: TargetServer:8080

### 3. Running the Rate Limiting Proxy

Build and run the rate limiter container:

```bash
cd rateLimitedUtil
docker build -t ratelimiter .
docker run --name RateLimiter -p 8080:8080 ratelimiter
```
To test the rate limiting feature, on the client terminal, test with:

```bash
for ($i=1; $i -le 7; $i++) { echo "http://google.com" | go run Client.go ; Start-Sleep -Milliseconds 100 }
```

### 4. Running the Malicious Blocking Proxy

Build and run the malicious URL blocker container:

```bash
cd safeBrowseUtil
docker build -t maliciousblocker .
docker run --name MaliciousBlocker -p 8080:8080 maliciousblocker
```
To test the malicious blocking feature, on the client terminal, enter the URL: http://testsafebrowsing.appspot.com/s/malware.html

---

### Conclusion

This project demonstrates the implementation of a custom HTTP proxy server in Go, enhanced with Google Safe Browsing integration to detect and block malicious URLs in real time. By intercepting client requests, validating target destinations, and conditionally forwarding or rejecting traffic, the proxy acts as a lightweight security layer for HTTP-based communication.
Through this project, key concepts such as HTTP request handling, proxy-aware client configuration, Safe Browsing API usage, and containerized deployment with Docker were explored and applied. The result is a modular, testable, and extensible proxy system that lays the groundwork for more advanced features like HTTPS support, rate limiting, and traffic analytics.
This proxy serves as both a practical tool and a learning platform for backend development, network security, and system design.
