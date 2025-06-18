package healthchecker

import (
	"fmt"
	"log"
	"net/http"
	"time"

	serverPool "github.com/x-sushant-x/Balancer/pool"
	"github.com/x-sushant-x/Balancer/types"
)

type HealthChecker struct {
	interval time.Duration
	pool     *serverPool.ServerPool
}

func NewHealthChecker(interval time.Duration, pool *serverPool.ServerPool) *HealthChecker {
	return &HealthChecker{
		interval: interval,
		pool:     pool,
	}
}

func (hc *HealthChecker) CheckServersHealth() {
	servers := hc.pool.GetAllServers()

	for {
		log.Println("Starting Health Check")

		for _, server := range servers {
			if server.HealthCheckURL != "" {
				doHTTPRequest(server)
			} else {
				// Some kind of notification system
				log.Printf("Health Check URL not specified for server: %s\n", server.URL)
			}
		}

		fmt.Println()
		fmt.Println()
		fmt.Println("Server Status")

		time.Sleep(hc.interval)
	}

}

func doHTTPRequest(server *types.Server) {
	client := &http.Client{}

	req, err := http.NewRequest(http.MethodGet, server.HealthCheckURL, nil)
	if err != nil {
		// Some kind of notification system
		updateServerUnhealthyStatus(server)
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		// Some kind of notification system
		updateServerUnhealthyStatus(server)
		return
	}

	if resp.StatusCode == http.StatusOK && !server.IsHealthy {
		updateServerHealthyStatus(server)
	}

	if resp.StatusCode != http.StatusOK {
		updateServerUnhealthyStatus(server)
	}
}

func updateServerUnhealthyStatus(server *types.Server) {
	log.Printf("Health Check Failed: Server did not responded with 200 status code. Server: %v\n", server.HealthCheckURL)

	if server.IsHealthy {
		server.FailureCount++
	}

	if server.FailureCount >= server.UnhealthyAfter && server.IsHealthy {
		server.IsHealthy = false
	}
}

func updateServerHealthyStatus(server *types.Server) {
	server.SuccessCount++

	if server.SuccessCount >= server.HealthyAfter && !server.IsHealthy {
		server.IsHealthy = true
		server.SuccessCount = 0
		server.FailureCount = 0
	}
}
