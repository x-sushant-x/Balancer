package healthchecker

import (
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
		go doHTTPRequest(*server)
	}

}

func doHTTPRequest(server types.Server) (*http.Response, error) {
	client := &http.Client{}

	req, err := http.NewRequest(http.MethodGet, server.HealthCheckURL, nil)
	if err != nil {
		return nil, err
	}

	return client.Do(req)
}
