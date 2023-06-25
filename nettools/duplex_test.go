package nettools_test

import (
	"bytes"
	"context"
	"net"
	"testing"

	"golang.org/x/sync/errgroup"

	"github.com/abferm/candi/nettools"
)

func TestDuplex(t *testing.T) {
	aIn, aOut := net.Pipe()
	bIn, bOut := net.Pipe()

	// instanciate a duplexed connection that reads from aOut and writes to bIn
	duplex := nettools.Duplex(aOut, bIn)

	readExpected := []byte("read test")
	writeExpected := []byte("write test")

	readActual := make([]byte, len(readExpected))
	writeActual := make([]byte, len(writeExpected))

	eg, _ := errgroup.WithContext(context.Background())

	eg.Go(func() error {
		aIn.Write(readExpected)
		bOut.Read(writeActual)
		return nil
	})

	n, err := duplex.Read(readActual)
	if err != nil {
		t.Fail()
		t.Logf("Unexpected error on read: %s", err)
	}
	if n != len(readExpected) {
		t.Fail()
		t.Logf("Only read %d of %d expected bytes", n, len(readExpected))
	}

	n, err = duplex.Write(writeExpected)
	if err != nil {
		t.Fail()
		t.Logf("Unexpected error on write: %s", err)
	}
	if n != len(writeExpected) {
		t.Fail()
		t.Logf("Only wrote %d of %d expected bytes", n, len(writeExpected))
	}

	eg.Wait()

	if !bytes.Equal(readActual, readExpected) {
		t.Fail()
		t.Logf("Read bytes do not match %s != %s", string(readActual), string(readExpected))
	}
	if !bytes.Equal(writeActual, writeExpected) {
		t.Fail()
		t.Logf("Wtite bytes do not match %s != %s", string(readActual), string(readExpected))
	}
}
