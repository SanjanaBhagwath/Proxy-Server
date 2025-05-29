package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {
	var targetURL string
	fmt.Println("Please enter the target URL you wish to connect to")
	fmt.Scan(&targetURL)
	proxyURL := "http://localhost:9090"

	method := "GET"
	request, err1 := http.NewRequestWithContext(context.TODO(), method, proxyURL, nil)
	request.Host = targetURL
	fmt.Println("Final Request URL:", request.URL)
	fmt.Println("Final Host:", request.Host)
	if err1 != nil {
		log.Fatalf("Error in creating a request: %v", err1)
	}

	response, err2 := new(http.Client).Do(request)
	if err2 != nil {
		log.Fatalf("Error in sending the request: %v", err2)
	}
	fmt.Println("Request has been sent to the proxy server successfully")

	responseBody, err3 := io.ReadAll(response.Body)
	if response.StatusCode != http.StatusOK {
		fmt.Println("The response returned with status code: ", response.StatusCode)
	}
	if err3 != nil {
		log.Fatalf("Error in printing out the response: %v", err3)
	}

	fmt.Println("Response body length: ", len(responseBody))
	fmt.Println("Response received from the server: ", string(responseBody))
	fmt.Println("Response headers: ", response.Header)

}
