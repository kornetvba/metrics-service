package config

import (
	"errors"
	"flag"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type NetAddr struct {
	Host string
	Port int
}

func (a *NetAddr) String() string {
	return fmt.Sprintf("%s:%d", a.Host, a.Port)
}

func (a *NetAddr) Set(addr string) error {
	addrSpl := strings.Split(addr, ":")
	if len(addrSpl) != 2 {
		return errors.New("address is not valid")
	}
	port, err := strconv.Atoi(addrSpl[1])
	if err != nil {
		return err
	}

	a.Host = addrSpl[0]
	a.Port = port
	return nil
}

var Addr = &NetAddr{
	Host: "localhost",
	Port: 8080,
}
var ReportInterval = 10 * time.Second
var PollInterval = 2 * time.Second

func ParseFlagAgent() {
	flag.Var(Addr, "a", "localhost:8080")
	flag.DurationVar(&ReportInterval, "r", ReportInterval, "report interval")
	flag.DurationVar(&PollInterval, "p", PollInterval, "poll interval")
	flag.Parse()
}
func ParseFlagServer() {
	flag.Var(Addr, "a", "localhost:8080")
	flag.Parse()
}
