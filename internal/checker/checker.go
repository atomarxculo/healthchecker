package checker

import (
	"context"
	"time"
)

type Result struct {
	Name string
	Type string
	Target string
	Healthy bool
	StatusCode int
	Latency time.Duration
	Error string
}

type Checker interface {
	Check(ctx context.Context) *Result
	Name() string
}