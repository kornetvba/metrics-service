package main

import (
	"github.com/kornetvba/metrics-service/internal/agent"
	"log"
	"time"
)

const (
	pollInterval   time.Duration = 2
	reportInterval time.Duration = 10
)

func main() {
	// Канал для синхронизации
	go func() {
		agent.CollectMetrics(pollInterval)

	}()
	err := agent.ClientMetric(reportInterval)
	if err != nil {
		log.Fatal(err)
	}
}
