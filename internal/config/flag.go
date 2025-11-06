package config

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
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
