package rc522

import (
	"time"
	"github.com/op0xA5/rc522/hw"
)

type deadline struct {
	t time.Time
}
func newDeadline(d time.Duration) deadline {
	if d > 0 {
		return deadline{ time.Now().Add(d) }
	}
	return deadline{}
}
func (d deadline) SinceNow() time.Duration {
	if (d.t.IsZero()) {		
		return 0
	}
	return time.Until(d.t)
}
func (d deadline) Timeout() bool {
	return d.SinceNow() < 0
}

func waitRegister(dev DeviceIO, addr uint8, val, mask uint8, timeout time.Duration) error {
	deadline := newDeadline(timeout)

	for {
		got, err := dev.ReadRegister(addr)
		if err != nil {
			return err
		}
		if got & mask == val {
			return nil
		}
		if deadline.Timeout() {
			return timeoutErr
		}
		time.Sleep(50 * time.Millisecond)
	}
}

func waitRegisterSet(dev DeviceIO, addr uint8, val uint8, timeout time.Duration) error {
	deadline := newDeadline(timeout)

	for {
		got, err := dev.ReadRegister(addr)
		if err != nil {
			return err
		}
		if got & val > 0 {
			return nil
		}
		if deadline.Timeout() {
			return timeoutErr
		}
		time.Sleep(50 * time.Millisecond)
	}
}

func waitIdle(dev DeviceIO, timeout time.Duration) error {
	return waitRegister(dev, hw.CommandReg, hw.CmdIdle, hw.CommandReg, timeout)
}