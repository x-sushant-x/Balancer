package healthchecker

import (
	"time"

	"github.com/x-sushant-x/Balancer/types"
)

type HealthChecker struct {
	Interval       time.Duration
	Timeout        time.Duration
	RetryCount     int
	HealthyAfter   int
	UnhealthyAfter int
	Servers        []*types.Server
}
