package isotp

import (
	"fmt"
	"net"

	"github.com/abferm/candi"
	"golang.org/x/sys/unix"
)

type Bus struct {
	iface *net.Interface
}

func BusByName(name string) (*Bus, error) {
	iface, err := net.InterfaceByName(name)
	if err != nil {
		return nil, err
	}
	return &Bus{
		iface: iface,
	}, nil
}

func (bus Bus) Dial(addr net.Addr) (net.Conn, error) {
	if addr.Network() != network_isotp {
		return nil, fmt.Errorf("address specifies wrong netowrk type: %q", addr.Network())
	}
	tpAddr, ok := addr.(Addr)
	if !ok {
		return nil, fmt.Errorf("address must be of type isotp.Addr")
	}
	fd, err := unix.Socket(candi.PF_CAN, unix.SOCK_DGRAM, unix.CAN_ISOTP)
	if err != nil {
		return nil, err
	}

	// TODO: allow options to be passed in
	opts := DefaultOptions()
	err = SetSockoptCanIsoTpOpst(fd, unix.SOL_CAN_BASE+unix.CAN_ISOTP, sockoptCanIsotpOpts, opts)
	if err != nil {
		unix.Close(fd)
		return nil, fmt.Errorf("set opts failed: %w", err)
	}

	socAddr := tpAddr.SocketAddr(bus.iface.Index)

	err = unix.Bind(fd, socAddr)
	if err != nil {
		unix.Close(fd)
		return nil, fmt.Errorf("bind failed: %w", err)
	}

	return &conn{
		bus:  &bus,
		addr: tpAddr,
		fd:   fd,
	}, nil
}
