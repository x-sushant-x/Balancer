package balancer

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

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

	client := &http.Client{
		Timeout: time.Second * 30,
	}

	res, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to reach backend server", http.StatusBadGateway)
		return
	}

	byteResp, err := io.ReadAll(res.Body)
	if err != nil {
		http.Error(w, "Failed to read response body", http.StatusInternalServerError)
		return
	}

	defer res.Body.Close()

	for k, values := range res.Header {
		for _, value := range values {
			w.Header().Add(k, value)
		}
	}

	w.WriteHeader(res.StatusCode)

	fmt.Fprint(w, string(byteResp))
}
