package main

import (
	"github.com/kornetvba/metrics-service/internal/agent"
	"github.com/kornetvba/metrics-service/internal/config"
	"log"
)

func main() {
	config.ParseFlagAgent()
	// Канал для синхронизации
	go func() {
		agent.CollectMetrics(config.PollInterval)

	}()
	err := agent.ClientMetric(config.ReportInterval)
	if err != nil {
		log.Fatal(err)
	}
}
