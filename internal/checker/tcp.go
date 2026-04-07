package checker

import (
	"context"
	"fmt"
	"net"
	"time"
)

type TCPChecker struct {
	name string
	host string
	port string
	timeout time.Duration
}

func NewTCPChecker(name, hostPort string, timeout time.Duration) *TCPChecker {
	return &TCPChecker{
		name: name,
		host: hostPort,
		timeout: timeout,
	}
}

func (t *TCPChecker) Name() string  {
	return t.name
}

func (t *TCPChecker) Check(ctx context.Context) *Result  {
	start := time.Now()
	result := &Result{
		Name: t.name,
		Type: "tcp",
		Target: t.host,
		Healthy: false,
	}

	ctx, cancel := context.WithTimeout(ctx, t.timeout)
	defer cancel()

	var dialer net.Dialer
	conn, err := dialer.DialContext(ctx, "tcp", t.host)
	result.Latency = time.Since(start)

	if err != nil {
		result.Error = fmt.Sprintf("connection failed: %v", err)
		return result
	}

	defer conn.Close()

	result.Healthy = true
	return result

}