package isotp

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"
)

func TestCloseCancelsRead(t *testing.T) {
	bus, err := BusByName("vcan0")
	if err != nil {
		t.Skip("device vcan0 not available")
	}

	addr := NewAddr(1, 2)

	conn, err := bus.Dial(addr)
	if err != nil {
		t.Fatalf("unexpected error when dialing: %s", err)
	}

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
		defer cancel()
		<-ctx.Done()
		conn.Close()
	}()
	b := make([]byte, 5)
	t.Log("starting read")
	_, err = conn.Read(b)
	if !errors.Is(err, os.ErrClosed) {
		t.Fatalf("unexpected error when reading: %s", err)
	}
}
