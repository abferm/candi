package nettools

import (
	"errors"
	"net"
	"os"
	"time"
)

func NewFileConn(f *os.File, localAddr, remoteAddr net.Addr) *FileConn {
	return &FileConn{
		f:          f,
		localAddr:  localAddr,
		remoteAddr: remoteAddr,
	}
}

type FileConn struct {
	localAddr, remoteAddr net.Addr
	f                     *os.File
}

// type enforcement
var _ net.Conn = new(FileConn)

// LocalAddr returns the local network address, if known.
func (c *FileConn) LocalAddr() net.Addr { return c.localAddr }

// RemoteAddr returns the remote network address, if known.
func (c *FileConn) RemoteAddr() net.Addr { return c.remoteAddr }

func (c *FileConn) Read(b []byte) (int, error) {
	n, err := c.f.Read(b)
	return n, AsNetOpError(err, c, "read")
}

func (c *FileConn) Write(b []byte) (int, error) {
	n, err := c.f.Write(b)
	return n, AsNetOpError(err, c, "write")
}

func (c *FileConn) SetDeadline(t time.Time) error {
	err := c.f.SetDeadline(t)
	return AsNetOpError(err, c, "set")
}

func (c *FileConn) SetReadDeadline(t time.Time) error {
	err := c.f.SetReadDeadline(t)
	return AsNetOpError(err, c, "set")
}

func (c *FileConn) SetWriteDeadline(t time.Time) error {
	err := c.f.SetWriteDeadline(t)
	return AsNetOpError(err, c, "set")
}

func (c *FileConn) Close() error {
	err := c.f.Close()
	return AsNetOpError(err, c, "close")
}

func AsNetOpError(err error, conn net.Conn, op string) *net.OpError {
	if err == nil {
		return nil
	}
	opErr := &net.OpError{
		Op:     op,
		Net:    conn.LocalAddr().Network(),
		Source: conn.LocalAddr(),
		Addr:   conn.RemoteAddr(),
		Err:    err,
	}
	// if this is a path error, unwrap the inner error and honor the Operation rather than the one provided
	var pe *os.PathError
	if errors.As(err, &pe) {
		opErr.Op = pe.Op
		opErr.Err = pe.Err
	}
	return opErr
}
