package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/abferm/candi/isotp"
)

func main() {
	err := Main()
	if err != nil {
		log.Fatalln(err)
	}
}

func Main() error {
	rx := flag.Uint("rxaddr", 0, "recieve address")
	tx := flag.Uint("txaddr", 0, "send address")
	iface := flag.String("interface", "vcan0", "can interface to use")
	flag.Parse()

	addr := isotp.NewAddr(uint32(*rx), uint32(*tx))
	bus, err := isotp.BusByName(*iface)
	if err != nil {
		return fmt.Errorf("failed to get bus: %w", err)
	}

	conn, err := bus.Dial(addr)
	if err != nil {
		return fmt.Errorf("dial failed: %w", err)
	}

	defer conn.Close()

	n, err := conn.Write([]byte{0xDE, 0xAD, 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0xBE, 0xEF})
	if err != nil {
		return fmt.Errorf("write failed after %d bytes: %w", n, err)
	}
	return nil
}
