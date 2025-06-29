# 🛡️ Go-Based Secure Proxy Server

This project is a modular, high-performance HTTP proxy server implemented in **Go**, designed to demonstrate backend concepts like request interception, IP masking, rate limiting, and malicious URL blocking using the **Google Safe Browsing API**. The system is fully containerized using **Docker**, with separate containers for each core component to simulate a scalable, real-world network.

## 🚀 Why Go?

Go (Golang) is well-suited for building networked services due to its:

- Lightweight concurrency model using goroutines
- Fast compilation and execution
- Rich standard library for HTTP, networking, and cryptography
- Simplicity and readability, making it ideal for scalable microservices

These features make Go a natural fit for implementing a performant and modular proxy server.

## 🧱 Architecture Overview

The system is composed of four Docker containers, each representing a distinct component in the proxy chain:

| Container           | Role                                                                 |
|---------------------|----------------------------------------------------------------------|
| `proxy-server`      | Main HTTP proxy that forwards requests and masks client IPs         |
| `target-server`     | Simulated backend server that responds to forwarded requests         |
| `rate-limiter`      | Middleware proxy that limits requests per client IP                 |
| `malicious-blocker` | Middleware proxy that blocks requests to malicious URLs using the Google Safe Browsing API |

All containers are connected within the same Docker network to simulate realistic routing and isolation.

