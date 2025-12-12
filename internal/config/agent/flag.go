package agent

import (
	"flag"
	"github.com/kornetvba/metrics-service/internal/config"
	"log"
	"os"
	"strconv"
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

	if addr, ok := os.LookupEnv("ADDRESS"); ok {
		err := AddrAgent.Set(addr)
		if err != nil {
			log.Print(err)
		}
	}
	if repInt, ok := os.LookupEnv("REPORT_INTERVAL"); ok {
		rep, err := strconv.Atoi(repInt)
		if err != nil {
			log.Print(err)
		}
		ReportInterval = rep
	}
	if pollInt, ok := os.LookupEnv("POLL_INTERVAL"); ok {
		poll, err := strconv.Atoi(pollInt)
		if err != nil {
			log.Print(poll)
		}
		PollInterval = poll
	}

}
