package checker

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type HTTPChecker struct {
	name string
	url string
	method string
	timeout time.Duration
}

func NewHTTPChecker(name, url, method string, timeout time.Duration) *HTTPChecker {
	if method == "" {
		method = "GET"
	}

	return &HTTPChecker{
		name: name,
		url: url,
		method: method,
		timeout: timeout,
	}
}

func (h *HTTPChecker) Name() string  {
	return h.name
}

func (h *HTTPChecker) Check(ctx context.Context) *Result  {
	start := time.Now()
	result := &Result{
		Name: h.name,
		Type: "http",
		Target: h.url,
		Healthy: false,
	}

	ctx, cancel := context.WithTimeout(ctx, h.timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, h.method, h.url, nil)
	if err != nil{
		result.Latency = time.Since(start)
		result.Error = fmt.Sprintf("creating request: %v", err)
		return result
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	result.Latency = time.Since(start)

	if err != nil {
		result.Error = fmt.Sprintf("creating request: %v", err)
		return result
	}
	defer resp.Body.Close()

	result.StatusCode = resp.StatusCode
	result.Healthy = resp.StatusCode >= 200 && resp.StatusCode < 400

	if !result.Healthy {
		result.Error = fmt.Sprintf("unhealthy status code: %d", resp.StatusCode)
	}

	return result

}