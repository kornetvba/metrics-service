package agent

import (
	"flag"
	"github.com/kornetvba/metrics-service/internal/config"
)

var AddrAgent = &config.NetAddr{
	Host: "localhost",
	Port: 8080,
}

var ReportInterval = 10
var PollInterval = 2

func ParseFlagAgent() {
	flag.Var(AddrAgent, "a", "localhost:8080")
	flag.IntVar(&ReportInterval, "r", ReportInterval, "report interval")
	flag.IntVar(&PollInterval, "p", PollInterval, "report interval")
	flag.Parse()
}
