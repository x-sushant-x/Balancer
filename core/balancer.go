package balancer

import (
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/x-sushant-x/Balancer/types"
)

type BalancerStrategy interface {
	GetNextServer() (*types.Server, error)
}

type LoadBalancer struct {
	strategy BalancerStrategy
}

func NewLoadBalancer(strategy BalancerStrategy) LoadBalancer {
	return LoadBalancer{
		strategy: strategy,
	}
}

func (lb *LoadBalancer) Serve(w http.ResponseWriter, r *http.Request) {
	server, err := lb.strategy.GetNextServer()
	if err != nil {
		http.Error(w, "internal server error", 500)
		return
	}

	targetURL, err := url.Parse(server.URL)
	if err != nil {
		http.Error(w, "Invalid backend URL", http.StatusInternalServerError)
		return
	}

	targetPath := strings.TrimRight(targetURL.String(), "/") + r.URL.Path
	if r.URL.RawQuery != "" {
		targetPath += "?" + r.URL.RawQuery
	}

	req, err := http.NewRequest(r.Method, targetPath, r.Body)
	if err != nil {
		http.Error(w, "Failed to create request to backend", http.StatusInternalServerError)
		return
	}

	for k, values := range r.Header {
		for _, value := range values {
			req.Header.Add(k, value)
		}
	}

	req.Header.Set("X-Forwarded-For", r.RemoteAddr)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, "Failed to reach backend server", http.StatusBadGateway)
		return
	}

	defer res.Body.Close()

	for k, values := range res.Header {
		for _, value := range values {
			w.Header().Add(k, value)
		}
	}

	w.WriteHeader(res.StatusCode)

	io.Copy(w, r.Body)
}
