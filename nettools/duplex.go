package nettools

import (
	"errors"
	"fmt"
	"net"
	"time"
)

type DuplexAddr struct {
	Read, Write net.Addr
}

func (a *DuplexAddr) Network() string {
	if a.Read.Network() == a.Write.Network() {
		return a.Read.Network()
	} else {
		return "read:" + a.Read.Network() + "+" + "write:" + a.Write.Network()
	}
}

func (a *DuplexAddr) String() string {
	return fmt.Sprintf("duplex{read:%s, write:%s}", a.Read, a.Write)
}

// Duplex pairs a set of connections together and uses one for read
// opperations and the other for write operations.
// This is useful for bi-directional communication with devices that
// do not support read/write with the same CAN ID pair.
func Duplex(read, write net.Conn) net.Conn {
	return &duplex{
		read:  read,
		write: write,
	}
}

type duplex struct {
	read, write net.Conn
}

// Read reads data from the duplexection.
// Read can be made to time out and return an error after a fixed
// time limit; see SetDeadline and SetReadDeadline.
func (c *duplex) Read(b []byte) (n int, err error) {
	return c.read.Read(b)
}

// Write writes data to the duplexection.
// Write can be made to time out and return an error after a fixed
// time limit; see SetDeadline and SetWriteDeadline.
func (c *duplex) Write(b []byte) (n int, err error) {
	return c.write.Write(b)
}

// Close closes the duplexection.
// Any blocked Read or Write operations will be unblocked and return errors.
func (c *duplex) Close() error {
	return errors.Join(c.read.Close(), c.write.Close())
}

// LocalAddr returns the local network address, if known.
func (c *duplex) LocalAddr() net.Addr {
	return &DuplexAddr{
		Read:  c.read.LocalAddr(),
		Write: c.write.LocalAddr(),
	}
}

// RemoteAddr returns the remote network address, if known.
func (c *duplex) RemoteAddr() net.Addr {
	return &DuplexAddr{
		// note crossover
		Write: c.read.RemoteAddr(),
		Read:  c.write.RemoteAddr(),
	}
}

// SetDeadline sets the read and write deadlines associated
// with the duplexection. It is equivalent to calling both
// SetReadDeadline and SetWriteDeadline.
//
// A deadline is an absolute time after which I/O operations
// fail instead of blocking. The deadline applies to all future
// and pending I/O, not just the immediately following call to
// Read or Write. After a deadline has been exceeded, the
// duplexection can be refreshed by setting a deadline in the future.
//
// If the deadline is exceeded a call to Read or Write or to other
// I/O methods will return an error that wraps os.ErrDeadlineExceeded.
// This can be tested using errors.Is(err, os.ErrDeadlineExceeded).
// The error's Timeout method will return true, but note that there
// are other possible errors for which the Timeout method will
// return true even if the deadline has not been exceeded.
//
// An idle timeout can be implemented by repeatedly extending
// the deadline after successful Read or Write calls.
//
// A zero value for t means I/O operations will not time out.
func (c *duplex) SetDeadline(t time.Time) error {
	return errors.Join(
		c.SetReadDeadline(t),
		c.SetWriteDeadline(t),
	)
}

// SetReadDeadline sets the deadline for future Read calls
// and any currently-blocked Read call.
// A zero value for t means Read will not time out.
func (c *duplex) SetReadDeadline(t time.Time) error {
	return c.read.SetReadDeadline(t)
}

// SetWriteDeadline sets the deadline for future Write calls
// and any currently-blocked Write call.
// Even if write times out, it may return n > 0, indicating that
// some of the data was successfully written.
// A zero value for t means Write will not time out.
func (c *duplex) SetWriteDeadline(t time.Time) error {
	return c.write.SetWriteDeadline(t)
}
