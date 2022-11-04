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
	buff := make([]byte, 12)
	n, err := conn.Read(buff)
	if err != nil {
		return fmt.Errorf("write failed after %d bytes: %w", n, err)
	}
	fmt.Printf("Received: %x\n", buff)
	return nil
}
