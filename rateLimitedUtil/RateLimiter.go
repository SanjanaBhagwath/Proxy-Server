package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"
)

// IPRequestState holds the request count and last request time for an IP.
type IPRequestState struct {
	Count      int
	LastAccess time.Time
}

// RateLimiter manages the state for basic rate limiting.
type RateLimiter struct {
	mu           sync.Mutex
	ipStates     map[string]*IPRequestState
	requestLimit int           // Max requests per time window
	window       time.Duration // Time window for the limit
}

// NewRateLimiter creates and initializes a new RateLimiter.
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		ipStates:     make(map[string]*IPRequestState),
		requestLimit: limit,
		window:       window,
	}
}

// CheckAndRecord checks if an IP is exceeding the rate limit and records the current request.
// Returns true if the request is allowed, false if it should be blocked.
func (rl *RateLimiter) CheckAndRecord(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	state, exists := rl.ipStates[ip]
	if !exists {
		state = &IPRequestState{
			Count:      0,
			LastAccess: time.Now(),
		}
		rl.ipStates[ip] = state
	}

	currentTime := time.Now()

	// If the last access was outside the current window, reset the count.
	if currentTime.Sub(state.LastAccess) > rl.window {
		state.Count = 0
	}

	state.Count++
	state.LastAccess = currentTime // Update last access time for the current request

	// Check if the request count exceeds the limit within the window.
	if state.Count > rl.requestLimit {
		fmt.Printf("Blocked IP %s: Too many requests (%d) within %s\n", ip, state.Count, rl.window)
		return false // Block the request
	}

	// Clean up old IP states to prevent unbounded memory growth.
	// This simple cleanup iterates on every request. For very high throughput,
	// consider moving this to a separate goroutine that runs periodically.
	for ipAddr, s := range rl.ipStates {
		if currentTime.Sub(s.LastAccess) > 5*rl.window { // If no activity for 5 windows, clean up
			delete(rl.ipStates, ipAddr)
		}
	}

	return true // Request is allowed
}

type proxyHandler struct {
	limiter *RateLimiter // Add a pointer to our RateLimiter
}

func (ph *proxyHandler) ServeHTTP(w http.ResponseWriter, requestReceived *http.Request) {
	clientIP, _, err := net.SplitHostPort(requestReceived.RemoteAddr)
	if err != nil {
		fmt.Printf("Could not parse client IP from RemoteAddr: %v. Using full RemoteAddr.\n", err)
		clientIP = requestReceived.RemoteAddr
	}

	// Check if the request is allowed by the rate limiter
	if ph.limiter != nil && !ph.limiter.CheckAndRecord(clientIP) {
		http.Error(w, "Too many requests.", http.StatusTooManyRequests)
		fmt.Printf("Request from %s blocked by rate limiter.\n", clientIP)
		return // Stop processing this request
	}

	fmt.Println("Request has been intercepted by the proxy server")
	fmt.Println("Final host from the request received: ", requestReceived.URL.String())
	fmt.Println("Request is received from the IP: ", clientIP)

	hostHeader := requestReceived.Host

	target := &url.URL{
		Scheme:   "http",
		Host:     hostHeader,
		Path:     requestReceived.URL.Path,
		RawQuery: requestReceived.URL.RawQuery,
	}

	targetURL := target.String()

	fmt.Println("Final target URL: ", targetURL)
	fmt.Println("Request received is: ", requestReceived)

	proxyRequest, err1 := http.NewRequestWithContext(context.TODO(), "GET", targetURL, nil)
	if err1 != nil {
		log.Fatalf("Error in creating the request to be forwarded: %v", err1)
	}

	// Copy headers from original request to proxy request
	for name, values := range requestReceived.Header {
		for _, value := range values {
			proxyRequest.Header.Add(name, value)
		}
	}
	// Important: Set Host header correctly for the upstream server
	proxyRequest.Host = hostHeader
	proxyRequest.URL.Host = hostHeader
	proxyRequest.URL.Scheme = "http" // Or "https" if you want to support HTTPS upstream

	// If the original request has a body (e.g., POST, PUT), copy it.
	if requestReceived.Body != nil {
		proxyRequest.Body = requestReceived.Body
	}

	response, err2 := new(http.Client).Do(proxyRequest)
	if err2 != nil {
		// Log the actual error for debugging
		fmt.Printf("Error in sending the request to target %s: %v\n", targetURL, err2)
		http.Error(w, "Proxy error: Could not reach target server.", http.StatusBadGateway)
		return
	}
	fmt.Println("Request has been sent to the server successfully")
	fmt.Println("Response has been received by the proxy")

	//Relaying responses back to the Client
	if response.StatusCode != http.StatusOK {
		fmt.Printf("The response returned with status code: %d from target %s\n", response.StatusCode, targetURL)
	}

	//Set the status code
	w.WriteHeader(response.StatusCode)

	//Set the headers
	for key, values := range response.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	//Copy the body
	_, copyErr := io.Copy(w, response.Body)
	if copyErr != nil {
		fmt.Printf("Error copying response body: %v\n", copyErr)
	}

	fmt.Println("Response has been relayed back to the client")

	response.Body.Close()
}

func main() {
	fmt.Println("Creating a listener on localhost at port 9090")

	// --- Rate Limiting Configuration ---
	requestLimit := 5              // Allow 5 requests
	timeWindow := 10 * time.Second // within a 10-second window per IP

	limiter := NewRateLimiter(requestLimit, timeWindow)

	proxyServer := http.Server{
		Addr:    "0.0.0.0:9090",
		Handler: &proxyHandler{limiter: limiter}, // Pass the limiter to the handler
	}

	fmt.Printf("Proxy server starting on %s.\n", proxyServer.Addr)
	fmt.Printf("Configured Rate Limit: %d requests per %v.\n", requestLimit, timeWindow)

	// Use log.Fatal for ListenAndServe to ensure the error is reported and program exits
	log.Fatal(proxyServer.ListenAndServe())
}
