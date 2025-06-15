package core

import (
	"errors"
	"sync"

	pool "github.com/x-sushant-x/Balancer/pool"
	"github.com/x-sushant-x/Balancer/types"
)

var (
	idx        = 0
	once       = sync.Once{}
	RoundRobin RoundRobinBalancer
)

type RoundRobinBalancer struct {
	pool *pool.ServerPool
	mu   sync.Mutex
}

func NewRoundRobinBalancer(pool *pool.ServerPool) {
	once.Do(func() {
		RoundRobin = RoundRobinBalancer{
			pool: pool,
		}
	})
}

func (rb *RoundRobinBalancer) GetNextServer() (*types.Server, error) {
	servers := rb.pool.GetAllServers()

	if len(servers) == 0 {
		return nil, errors.New("no servers found")
	}

	rb.mu.Lock()
	idx = (idx + 1) % len(servers)
	defer rb.mu.Unlock()

	selectedServer := servers[idx]

	return selectedServer, nil
}
