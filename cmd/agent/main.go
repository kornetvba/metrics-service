package main

import (
	"github.com/kornetvba/metrics-service/internal/agent"
	"log"
	"time"
)

const (
	PollInterval   time.Duration = 2
	ReportInterval time.Duration = 10
)

func main() {
	// Канал для синхронизации
	go func() {
		agent.CollectMetrics(PollInterval)

	}()
	err := agent.ClientMetric(ReportInterval)
	if err != nil {
		log.Fatal(err)
	}
}
