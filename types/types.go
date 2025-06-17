package types

import "time"

type Server struct {
	ID              string
	Name            string
	Protocol        string
	Host            string
	Port            int
	URL             string
	IsHealthy       bool
	LastHealthCheck time.Time
	FailureCount    int
	SuccessCount    int
}
