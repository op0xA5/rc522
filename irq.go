package rc522

import (
	"time"
	"github.com/op0xA5/rc522/hw"
)

type irq int
const (
	irqTx      = irq(hw.ComIrqReg)<<8 | irq(hw.TxIRq)
	irqRx      = irq(hw.ComIrqReg)<<8 | irq(hw.RxIRq)
	irqIdle    = irq(hw.ComIrqReg)<<8 | irq(hw.IdleIRq)
	irqHiAlert = irq(hw.ComIrqReg)<<8 | irq(hw.HiAlertIRq)
	irqLoAlert = irq(hw.ComIrqReg)<<8 | irq(hw.LoAlertIRq)
	irqErr     = irq(hw.ComIrqReg)<<8 | irq(hw.ErrIRq)
	irqTimer   = irq(hw.ComIrqReg)<<8 | irq(hw.TimerIRq)
	irqMfinAct = irq(hw.DivIrqReg)<<8 | irq(hw.MfinActIRq)
	irqCRC     = irq(hw.DivIrqReg)<<8 | irq(hw.CRCIRq)

	irqCom  = irq(hw.ComIrqReg<<8) | irq(0x7F)
)
func (irq irq) addr() uint8 {
	return uint8(irq >> 8)
}
func (irq irq) val() uint8 {
	return uint8(irq)
}
func (irq irq) wait(dev DeviceIO, timeout time.Duration) (error) {
	return waitRegisterSet(dev, irq.addr(), irq.val(), timeout)
}
func (irq irq) clear(dev DeviceIO) error {
	return dev.WriteRegister(irq.addr(), irq.val() & 0x7F)
}
