package main

import (
	"github.com/kornetvba/metrics-service/internal/agent"
	flag "github.com/kornetvba/metrics-service/internal/config/agent"
	"log"
)

func main() {
	flag.ParseFlagAgent()
	// Канал для синхронизации
	go func() {
		agent.CollectMetrics(flag.PollInterval)

	}()
	err := agent.ClientMetric(flag.ReportInterval)
	if err != nil {
		log.Fatal(err)
	}
}
