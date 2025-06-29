package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

type proxyHandler struct{}

func (ph *proxyHandler) ServeHTTP(w http.ResponseWriter, requestReceived *http.Request) {

	fmt.Println("Request has been intercepted by the proxy server")
	fmt.Println("Final host from the request received: ", requestReceived.URL.String())

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
	// Print IP address of the client from which request is received
	fmt.Println("Request is received from the IP: ", requestReceived.RemoteAddr)

	proxyRequest, err1 := http.NewRequestWithContext(context.TODO(), "GET", targetURL, nil)
	if err1 != nil {
		log.Fatalf("Error in creating the request to be forwarded: %v", err1)
	}

	proxyRequest.Host = hostHeader

	response, err2 := new(http.Client).Do(proxyRequest)
	if err2 != nil {
		log.Fatalf("Error in sending the request: %v", err2)
	}
	fmt.Println("Request has been sent to the server successfully")
	fmt.Println("Response has been received by the proxy")

	//Relaying responses back to the Client

	if response.StatusCode != http.StatusOK {
		fmt.Println("The response returned with status code: ", response.StatusCode)
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
	io.Copy(w, response.Body)

	fmt.Println("Response has been relayed back to the client")

	response.Body.Close()

	//fmt.Println("Response body length: ", len(responseBody))
	//fmt.Println("Response received from the server: ", string(responseBody))
	//fmt.Println("Response headers: ", response.Header)
}

func main() {
	fmt.Println("Creating a listener on localhost at port 9090")

	proxyServer := http.Server{
		Addr:    "0.0.0.0:9090",
		Handler: &proxyHandler{},
	}

	defer proxyServer.Close()
	go proxyServer.ListenAndServe()

	select {}

}
