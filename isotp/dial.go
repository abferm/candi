package isotp

import (
	"fmt"
	"net"
)

func Dial(busName string, addr net.Addr) (net.Conn, error) {
	bus, err := BusByName(busName)
	if err != nil {
		return nil, fmt.Errorf("error acquiring bus: %w", err)
	}
	return bus.Dial(addr)
}

func DialDuplex(busName string, read, write net.Addr) (net.Conn, error) {
	bus, err := BusByName(busName)
	if err != nil {
		return nil, fmt.Errorf("error acquiring bus: %w", err)
	}
	return bus.DialDuplex(read, write)
}
