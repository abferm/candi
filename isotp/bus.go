package isotp

import (
	"fmt"
	"net"
	"syscall"
	"unsafe"

	"github.com/abferm/candi"
	"golang.org/x/sys/unix"
)

type Bus struct {
	iface     *net.Interface
	localAddr Addr
}

func BusByName(name string, localAddr Addr) (*Bus, error) {
	iface, err := net.InterfaceByName(name)
	if err != nil {
		return nil, err
	}
	return &Bus{
		iface:     iface,
		localAddr: localAddr,
	}, nil
}

func (bus Bus) LocalAddr() net.Addr { return bus.localAddr }

func (bus Bus) Dial(addr net.Addr) (net.Conn, error) {
	if addr.Network() != network_isotp {
		return nil, fmt.Errorf("address specifies wrong netowrk type: %q", addr.Network())
	}
	txAddr, ok := addr.(isotpAddr)
	if !ok {
		return nil, fmt.Errorf("address must be of type isotp.Addr")
	}
	fd, err := unix.Socket(candi.PF_CAN, unix.SOCK_DGRAM, unix.CAN_ISOTP)
	if err != nil {
		return nil, err
	}

	opts := DefaultOptions()
	err = SetSockoptCanIsoTpOpst(fd, unix.SOL_CAN_BASE+unix.CAN_ISOTP, CAN_ISOTP_OPTS, opts)
	if err != nil {
		unix.Close(fd)
		return nil, fmt.Errorf("set opts failed: %w", err)
	}
	//TODO unix.Setsocopt( sol_can_isotp)
	// see https://github.com/linux-can/can-utils/blob/aa3f0299251953e4b3e9023cdaaf80ca8354718e/isotpsend.c#L251
	socAddr := &unix.SockaddrCAN{
		Ifindex: bus.iface.Index,
		RxID:    uint32(bus.localAddr),
		TxID:    uint32(txAddr.Addr()),
	}
	err = unix.Bind(fd, socAddr)
	if err != nil {
		unix.Close(fd)
		return nil, fmt.Errorf("bind failed: %w", err)
	}

	return &conn{
		bus:    &bus,
		txAddr: txAddr.Addr(),
		fd:     fd,
	}, nil
}

func SetSockoptCanIsoTpOpst(s int, level int, name int, val *Options) error {
	return setsockopt(s, level, name, unsafe.Pointer(val), unsafe.Sizeof(*val))
}

func setsockopt(s int, level int, name int, val unsafe.Pointer, vallen uintptr) (err error) {
	_, _, e1 := unix.Syscall6(unix.SYS_SETSOCKOPT, uintptr(s), uintptr(level), uintptr(name), uintptr(val), uintptr(vallen), 0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}

// errnoErr returns common boxed Errno values, to prevent
// allocations at runtime.
func errnoErr(e syscall.Errno) error {
	switch e {
	case 0:
		return nil
	case unix.EAGAIN:
		return syscall.EAGAIN
	case unix.EINVAL:
		return syscall.EINVAL
	case unix.ENOENT:
		return syscall.ENOENT
	}
	return e
}
