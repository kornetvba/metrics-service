package server

import (
	"flag"
	"github.com/kornetvba/metrics-service/internal/config"
	"log"
	"os"
)

var AddrServer = &config.NetAddr{
	Host: "localhost",
	Port: 8080,
}

func ParseFlagServer() {
	flag.Var(AddrServer, "a", "localhost:8080")
	flag.Parse()

	if addr, ok := os.LookupEnv("ADDRESS"); ok {
		err := AddrServer.Set(addr)
		if err != nil {
			log.Print(err)
		}
	}

}
