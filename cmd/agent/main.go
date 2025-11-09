package main

import (
	"github.com/kornetvba/metrics-service/internal/agent"
	flag "github.com/kornetvba/metrics-service/internal/config/agent"
	"log"
	"time"
)

func main() {
	flag.ParseFlagAgent()
	// Канал для синхронизации
	go func() {
		agent.CollectMetrics(time.Duration(flag.PollInterval) * time.Second)

	}()
	err := agent.ClientMetric(time.Duration(flag.ReportInterval) * time.Second)
	if err != nil {
		log.Fatal(err)
	}
}
