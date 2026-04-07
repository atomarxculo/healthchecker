package main

import (
	"context"
	"flag"
	"fmt"
	"healthchecker/internal/checker"
	"healthchecker/internal/config"
	"log"
	"os"
	"strings"
	"time"
)

func main()  {
	var (
		url string
		method string
		timeout time.Duration
		verbose bool
		configFile string
	)

	flag.StringVar(&url, "url", "", "URL to check (required)")
	flag.StringVar(&method, "method", "GET", "Method for the request")
	flag.DurationVar(&timeout, "timeout", 5*time.Second, "Timeout for the request")
	flag.BoolVar(&verbose, "verbose", false, "Enable verbose output")
	flag.StringVar(&configFile, "config", "", "Path to YAML config file")

	flag.Parse()

	if configFile != "" {
		runWithConfig(configFile, verbose)
		return
	}

	if url == "" {
		fmt.Fprintf(os.Stderr, "Error: either -url or -config must be provided\n")
		flag.Usage()
		os.Exit(1)
	}

	method = strings.ToUpper(method)
	if method != "GET" && method != "HEAD" {
		fmt.Fprintf(os.Stderr, "Error: method must be GET or HEAD\n")
		os.Exit(1)
	}

	httpchecker := checker.NewHTTPChecker("cli-check", url, method, timeout)
	
	ctx := context.Background()
	result := httpchecker.Check(ctx)

	if verbose {
		log.Printf("[VERBOSE] Response time: %v\n", result.Latency)
		if result.StatusCode > 0 {
			log.Printf("[VERBOSE] Status code: %d\n", result.StatusCode)
		}
	}

	if result.Healthy {
		fmt.Println("OK")
		os.Exit(0)
	} else {
		fmt.Fprintf(os.Stderr, "FAIL: %s\n", result.Error)
		os.Exit(1)
	}
}

func runWithConfig(configPath string, verbose bool)  {
	cfg, err := config.Load(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}
	if verbose {
		log.Printf("Loaded %d services from config", len(cfg.Services))
	}
	var checkers []checker.Checker

	for _, svc := range cfg.Services{
		timeout := svc.Timeout
		if timeout == 0 {
			timeout = cfg.Global.Timeout
			if timeout == 0 {
				timeout = 5 * time.Second
			}
		}

		switch svc.Type {
		case "http":
			checkers = append(checkers, checker.NewHTTPChecker(svc.Name, svc.Target, svc.Method, timeout))
		case "tcp":
			checkers = append(checkers, checker.NewTCPChecker(svc.Name, svc.Target, timeout))
		default:
			fmt.Fprintf(os.Stderr, "Warning: unknown service type %s for %s\n", svc.Type, svc.Name)
		}
	}

	ctx := context.Background()
	allHealthy := true

	for _, ch := range checkers {
		result := ch.Check(ctx)

		if verbose {
			log.Printf("[VERBOSE] Check %s: healthy=%v, latency=%v", result.Name, result.Healthy, result.Latency)
		}

		if result.Healthy {
			fmt.Printf("%s: OK\n", result.Name)
		} else {
			fmt.Printf("%s: FAIL (%s)\n", result.Name, result.Error)
			allHealthy = false
		}
	}
	if allHealthy {
		os.Exit(0)
	} else {
		os.Exit(1)
	}
}

