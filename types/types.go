package types

import "time"

type Server struct {
	ID              string        `json:"id"`
	Name            string        `json:"name"`
	Protocol        string        `json:"protocol"`
	Host            string        `json:"host"`
	Port            int           `json:"port"`
	URL             string        `json:"url"`
	HealthCheckURL  string        `json:"health_check_url"`
	IsHealthy       bool          `json:"is_healthy"`
	Timeout         time.Duration `json:"timeout"`
	LastHealthCheck time.Time     `json:"last_health_check"`
	FailureCount    int           `json:"failure_count"`
	SuccessCount    int           `json:"success_count"`
	HealthyAfter    int           `json:"healthy_after"`
	UnhealthyAfter  int           `json:"unhealthy_after"`
	RetryCount      int           `json:"retry_count"`
}
