package main

import (
	"fmt"
	"net/http"
)

type targetHandler struct{}

func (th *targetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Target Server received request from:", r.RemoteAddr) // Logs Proxy's IP

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Hello from Target Server!") // Sample response
}

func main() {
	fmt.Println("Creating a listener on localhost at port 8080")

	targetServer := http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: &targetHandler{},
	}

	defer targetServer.Close()
	go targetServer.ListenAndServe()

	select {}
}
