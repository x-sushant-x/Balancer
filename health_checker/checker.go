package healthchecker

import (
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

	for _, server := range servers {
		doHTTPRequest(server)
	}

}

func doHTTPRequest(server *types.Server) {
	client := &http.Client{}

	req, err := http.NewRequest(http.MethodGet, server.HealthCheckURL, nil)
	if err != nil {
		// Some kind of notification system
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		// Some kind of notification system
		return
	}

	if resp.StatusCode == http.StatusOK && !server.IsHealthy {
		server.SuccessCount++

		if server.SuccessCount >= server.HealthyAfter {
			server.IsHealthy = true
			server.SuccessCount = 0
			server.FailureCount = 0
		}
	}

	if resp.StatusCode != http.StatusOK {
		log.Println("Health Check Failed: Server did not responded with 200 status code.")
		server.FailureCount++

		if server.FailureCount >= server.UnhealthyAfter && server.IsHealthy {
			server.IsHealthy = false
		}
	}
}
