package isotp

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/unix"
)

const (
	sockoptCanIsotpOpts = 1

	/* flags for isotp behaviour */

	CAN_ISOTP_LISTEN_MODE   = 0x0001 /* listen only (do not send FC) */
	CAN_ISOTP_EXTEND_ADDR   = 0x0002 /* enable extended addressing */
	CAN_ISOTP_TX_PADDING    = 0x0004 /* enable CAN frame padding tx path */
	CAN_ISOTP_RX_PADDING    = 0x0008 /* enable CAN frame padding rx path */
	CAN_ISOTP_CHK_PAD_LEN   = 0x0010 /* check received CAN frame padding */
	CAN_ISOTP_CHK_PAD_DATA  = 0x0020 /* check received CAN frame padding */
	CAN_ISOTP_HALF_DUPLEX   = 0x0040 /* half duplex error state handling */
	CAN_ISOTP_FORCE_TXSTMIN = 0x0080 /* ignore stmin from received FC */
	CAN_ISOTP_FORCE_RXSTMIN = 0x0100 /* ignore CFs depending on rx stmin */
	CAN_ISOTP_RX_EXT_ADDR   = 0x0200 /* different rx extended addressing */
	CAN_ISOTP_WAIT_TX_DONE  = 0x0400 /* wait for tx completion */
	CAN_ISOTP_SF_BROADCAST  = 0x0800 /* 1-to-N functional addressing */
	CAN_ISOTP_CF_BROADCAST  = 0x1000 /* 1-to-N transmission w/o FC */

	/* protocol machine default values */

	CAN_ISOTP_DEFAULT_FLAGS        = CAN_ISOTP_WAIT_TX_DONE
	CAN_ISOTP_DEFAULT_EXT_ADDRESS  = 0x00
	CAN_ISOTP_DEFAULT_PAD_CONTENT  = 0xCC  /* prevent bit-stuffing */
	CAN_ISOTP_DEFAULT_FRAME_TXTIME = 50000 /* 50 micro seconds */
	CAN_ISOTP_DEFAULT_RECV_BS      = 0
	CAN_ISOTP_DEFAULT_RECV_STMIN   = 0x00
	CAN_ISOTP_DEFAULT_RECV_WFTMAX  = 0
)

type Options struct {
	flags uint32 /* set flags for isotp behaviour.	*/
	/* __u32 value : flags see below	*/

	frameTxTime uint32 /* frame transmission time (N_As/N_Ar)	*/
	/* __u32 value : time in nano secs	*/

	extAddress uint8 /* set address for extended addressing	*/
	/* __u8 value : extended address	*/

	txPadContent byte /* set content of padding byte (tx)	*/
	/* __u8 value : content	on tx path	*/

	rxPadContent byte /* set content of padding byte (rx)	*/
	/* __u8 value : content	on rx path	*/

	rxExtAddress uint8 /* set address for extended addressing	*/
	/* __u8 value : extended address (rx)	*/
}

func (o *Options) SetFlag(f uint32) {
	o.flags |= f
}

func (o *Options) ClearFlag(f uint32) {
	o.flags &= ^f
}

func (o *Options) SetTXPadding(content byte) {
	o.SetFlag(CAN_ISOTP_TX_PADDING)
	o.txPadContent = content
}

func (o *Options) SetRXPadding(content byte) {
	o.SetFlag(CAN_ISOTP_RX_PADDING)
	o.rxPadContent = content
}

func (o *Options) SetExtendedAddress(a uint8) {
	o.SetFlag(CAN_ISOTP_EXTEND_ADDR)
	o.extAddress = a
}

func (o *Options) SetExtendedRXAddress(a uint8) {
	o.SetFlag(CAN_ISOTP_RX_EXT_ADDR)
	o.rxExtAddress = a
}

func DefaultOptions() *Options {
	return &Options{
		flags:        CAN_ISOTP_DEFAULT_FLAGS,
		frameTxTime:  CAN_ISOTP_DEFAULT_FRAME_TXTIME,
		extAddress:   CAN_ISOTP_DEFAULT_EXT_ADDRESS,
		txPadContent: CAN_ISOTP_DEFAULT_PAD_CONTENT,
		rxPadContent: CAN_ISOTP_DEFAULT_PAD_CONTENT,
		rxExtAddress: CAN_ISOTP_DEFAULT_EXT_ADDRESS,
	}
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
