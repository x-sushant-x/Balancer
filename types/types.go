package types

import "time"

type Server struct {
	ID              string
	Name            string
	Host            string
	Port            int
	IsHealthy       bool
	LastHealthCheck time.Time
}
