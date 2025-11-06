package agent

import (
	"flag"
	"github.com/kornetvba/metrics-service/internal/config"
	"time"
)

var AddrAgent = &config.NetAddr{
	Host: "localhost",
	Port: 8080,
}

var ReportInterval = 10 * time.Second
var PollInterval = 2 * time.Second

func ParseFlagAgent() {
	flag.Var(AddrAgent, "a", "localhost:8080")
	flag.DurationVar(&ReportInterval, "r", ReportInterval, "report interval")
	flag.DurationVar(&PollInterval, "p", PollInterval, "poll interval")
	flag.Parse()
}
