package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	balancer "github.com/x-sushant-x/Balancer/core"
	healthchecker "github.com/x-sushant-x/Balancer/health_checker"
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
		HealthyAfter:    3,
		UnhealthyAfter:  3,
		HealthCheckURL:  "http://localhost:3001/health",
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
		HealthyAfter:    3,
		UnhealthyAfter:  3,
		HealthCheckURL:  "http://localhost:3002/health",
	}, {
		ID:              "3",
		Name:            "Server 3",
		Protocol:        "http",
		Host:            "localhost",
		Port:            3003,
		URL:             "http://localhost:3003",
		IsHealthy:       true,
		LastHealthCheck: time.Now(),
		HealthyAfter:    3,
		UnhealthyAfter:  3,
		HealthCheckURL:  "http://localhost:3003/health",
	},
}

func main() {
	pool := serverPool.NewServerPool()

	for _, server := range servers {
		pool.AddServer(&server)
	}

	rb := balancer.NewRoundRobinBalancer(pool)

	lb := balancer.NewLoadBalancer(rb)

	healthChecker := healthchecker.NewHealthChecker(time.Second*5, pool)
	go healthChecker.CheckServersHealth()

	distributeLoad(3000, lb)
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
