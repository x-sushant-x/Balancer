package core

import "github.com/x-sushant-x/Balancer/types"

type LoadBalancer interface {
	GetNextServer() (*types.Server, error)
}
