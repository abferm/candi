package isotp

import (
	"fmt"
	"net"
	"time"

	"golang.org/x/sys/unix"
)

type conn struct {
	bus           *Bus
	addr          Addr
	fd            int
	readDeadline  time.Time
	writeDeadline time.Time
}

// Read reads data from the connection.
// Read can be made to time out and return an error after a fixed
// time limit; see SetDeadline and SetReadDeadline.
func (c *conn) Read(b []byte) (n int, err error) {
	timeout := new(unix.Timeval)
	if c.readDeadline.IsZero() {
		timeout.Sec = 0
		timeout.Usec = 0
	} else if c.readDeadline.Before(time.Now()) {
		return 0, fmt.Errorf("deadline exceeded")
	} else {
		*timeout = unix.NsecToTimeval(time.Until(c.readDeadline).Nanoseconds())
	}
	err = unix.SetsockoptTimeval(c.fd, unix.SOL_SOCKET, unix.SO_RCVTIMEO, timeout)
	if err != nil {
		return 0, fmt.Errorf("failed to set read timeout: %w", err)
	}
	return unix.Read(c.fd, b)
}

// Write writes data to the connection.
// Write can be made to time out and return an error after a fixed
// time limit; see SetDeadline and SetWriteDeadline.
func (c *conn) Write(b []byte) (n int, err error) {
	timeout := new(unix.Timeval)
	if c.writeDeadline.IsZero() {
		timeout.Sec = 0
		timeout.Usec = 0
	} else if c.writeDeadline.Before(time.Now()) {
		return 0, fmt.Errorf("deadline exceeded")
	} else {
		*timeout = unix.NsecToTimeval(time.Until(c.writeDeadline).Nanoseconds())
	}
	err = unix.SetsockoptTimeval(c.fd, unix.SOL_SOCKET, unix.SO_RCVTIMEO, timeout)
	if err != nil {
		return 0, fmt.Errorf("failed to set write timeout: %w", err)
	}
	return unix.Write(c.fd, b)
}

// Close closes the connection.
// Any blocked Read or Write operations will be unblocked and return errors.
func (c *conn) Close() error {
	return unix.Close(c.fd)
}

// LocalAddr returns the local network address, if known.
func (c *conn) LocalAddr() net.Addr { return c.addr }

// RemoteAddr returns the remote network address, if known.
func (c *conn) RemoteAddr() net.Addr {
	// remote addr is the same as ours, but with swapped tx/rx ids
	return NewAddr(c.addr.TxID, c.addr.RxID)
}

// SetDeadline sets the read and write deadlines associated
// with the connection. It is equivalent to calling both
// SetReadDeadline and SetWriteDeadline.
//
// A deadline is an absolute time after which I/O operations
// fail instead of blocking. The deadline applies to all future
// and pending I/O, not just the immediately following call to
// Read or Write. After a deadline has been exceeded, the
// connection can be refreshed by setting a deadline in the future.
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
func (c *conn) SetDeadline(t time.Time) error {
	c.SetReadDeadline(t)
	c.SetWriteDeadline(t)
	return nil
}

// SetReadDeadline sets the deadline for future Read calls
// and any currently-blocked Read call.
// A zero value for t means Read will not time out.
func (c *conn) SetReadDeadline(t time.Time) error {
	c.readDeadline = t
	return nil
}

// SetWriteDeadline sets the deadline for future Write calls
// and any currently-blocked Write call.
// Even if write times out, it may return n > 0, indicating that
// some of the data was successfully written.
// A zero value for t means Write will not time out.
func (c *conn) SetWriteDeadline(t time.Time) error {
	c.writeDeadline = t
	return nil
}
