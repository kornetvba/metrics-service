package main

import (
	"context"
	"github.com/kornetvba/metrics-service/internal/agent"
	flag "github.com/kornetvba/metrics-service/internal/config/agent"
	"log"
	"time"
)

func main() {
	flag.ParseFlagAgent()
	// Канал для синхронизации
	go func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		agent.CollectMetrics(ctx, time.Duration(flag.PollInterval)*time.Second)

	}()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := agent.ClientMetric(ctx, time.Duration(flag.ReportInterval)*time.Second)
	if err != nil {
		log.Fatal(err)
	}
}
