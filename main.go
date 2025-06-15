package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	balancer "github.com/x-sushant-x/Balancer/core"
	serverPool "github.com/x-sushant-x/Balancer/pool"
	"github.com/x-sushant-x/Balancer/types"
)

var servers = []types.Server{
	{
		ID:              "1",
		Name:            "Server 1",
		Protocol:        "http",
		Host:            "localhost",
		Port:            3001,
		URL:             "http://localhost:3001",
		IsHealthy:       true,
		LastHealthCheck: time.Now(),
	},
	{
		ID:              "2",
		Name:            "Server 2",
		Protocol:        "http",
		Host:            "localhost",
		Port:            3002,
		URL:             "http://localhost:3002",
		IsHealthy:       true,
		LastHealthCheck: time.Now(),
	}, {
		ID:              "3",
		Name:            "Server 3",
		Protocol:        "http",
		Host:            "localhost",
		Port:            3003,
		URL:             "http://localhost:3003",
		IsHealthy:       true,
		LastHealthCheck: time.Now(),
	},
}

func main() {
	pool := serverPool.NewServerPool()

	for _, server := range servers {
		pool.AddServer(&server)
	}

	rb := balancer.NewRoundRobinBalancer(pool)

	lb := balancer.NewLoadBalancer(rb)
	distributeLoad(3000, lb)

	select {}
}

func distributeLoad(port int, lb balancer.LoadBalancer) {
	mux := http.NewServeMux()

	mux.HandleFunc("/", lb.Serve)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	log.Printf("Starting load balancer on port %d", port)

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Load balancer on port %d failed: %v", port, err)
	}
}
