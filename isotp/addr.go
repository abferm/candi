package isotp

import (
	"net"
	"strconv"
)

const network_isotp = "can_isotp"

type Addr uint32

func (a Addr) Network() string { return network_isotp }
func (a Addr) String() string  { return strconv.FormatUint(uint64(a), 10) }
func (a Addr) Addr() Addr      { return a }

type isotpAddr interface {
	net.Addr
	Addr() Addr
}
