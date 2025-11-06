package server

import (
	"flag"
	"github.com/kornetvba/metrics-service/internal/config"
)

var AddrServer = &config.NetAddr{
	Host: "localhost",
	Port: 8080,
}

func ParseFlagServer() {
	flag.Var(AddrServer, "a", "localhost:8080")
	flag.Parse()
}
