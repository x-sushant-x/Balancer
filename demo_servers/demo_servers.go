package main

import (
	"fmt"
	"log"
	"net/http"
)

func startServer(port int) {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello From Server With Port: %d", port)
	})

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	log.Printf("Starting server on port %d", port)

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server on port %d failed: %v", port, err)
	}
}

func main() {
	go startServer(3001)
	go startServer(3002)
	go startServer(3003)

	select {}
}
