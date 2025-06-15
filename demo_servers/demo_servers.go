package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

func startServer(port int) {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/v1", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		res := fmt.Sprintf("Hello From Server: %d", port)
		io.Copy(w, strings.NewReader(res))
	})

	mux.HandleFunc("/*", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		io.Copy(w, strings.NewReader("Not Found"))
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
