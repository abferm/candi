package isotp

import (
	"fmt"
	"net"

	"golang.org/x/sys/unix"
)

const network_isotp = "can_isotp"

type Addr struct {
	RxID uint32
	TxID uint32
}

func NewAddr(RxID, TxID uint32) Addr {
	return Addr{
		RxID: RxID,
		TxID: TxID,
	}
}

func (a Addr) Network() string { return network_isotp }
func (a Addr) String() string  { return fmt.Sprintf("%s{rx:%d,tx:%d", network_isotp, a.RxID, a.TxID) }
func (a Addr) SocketAddr(ifIndex int) *unix.SockaddrCAN {
	return &unix.SockaddrCAN{
		Ifindex: ifIndex,
		RxID:    a.RxID,
		TxID:    a.TxID,
	}
}

type isotpAddr interface {
	net.Addr
	SocketAddr(ifIndex int) *unix.SockaddrCAN
}
