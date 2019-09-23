package rc522

import (
	"fmt"
	"time"
	"github.com/op0xA5/rc522/hw"
)

const (
	cmdIdle             = hw.CmdIdle
	cmdMem              = hw.CmdMem
	cmdGenerateRandomID = hw.CmdGenerateRandomID
	cmdCalcCRC          = hw.CmdCalcCRC
	cmdTransmit         = hw.CmdTransmit
	cmdNoCmdChange      = hw.CmdNoCmdChange
	cmdReceive          = hw.CmdReceive
	cmdTransceive       = hw.CmdTransceive
	cmdMFAuthent        = hw.CmdMFAuthent
	cmdSoftReset        = hw.CmdSoftReset
)

func command(dev DeviceIO, cmd uint8) error {
	var commandNeedsRX bool
	switch (cmd) {
	case cmdIdle, cmdMem, cmdGenerateRandomID, cmdCalcCRC,
		cmdTransmit, cmdSoftReset:
		commandNeedsRX = false
	case cmdReceive, cmdTransceive, cmdMFAuthent:
		commandNeedsRX = true
	case cmdNoCmdChange:
		return nil
	default:
		return fmt.Errorf("unknown command: %x", cmd)
	}

	if commandNeedsRX == false {
		cmd |= hw.RcvOff
	}

	return dev.WriteRegister(hw.CommandReg, cmd)
}

/*
func CalculateCRC(dev DeviceIO, p []byte, timeout time.Duration) (uint16, error) {
	deadline := newDeadline(timeout)
	// stop any active command
	if err := command(dev, cmdIdle); err != nil {
		return 0, err
	}
	if err := irqCRC.clear(dev); err != nil {
		return 0, err
	}
	if err := flushFIFO(dev); err != nil {
		return 0, err
	}
	if err := command(dev, cmdCalcCRC); err != nil {
		return 0, err
	}
	if err := writeToFIFO(dev, p, deadline.SinceNow()); err != nil {
		return 0, err
	}
	if err := waitRegisterSet(dev, hw.Status1Reg, hw.CRCReady, deadline.SinceNow()); err != nil {
		return 0, err
	}
	//Stop calculating CRC for new content in the FIFO
	if err := command(dev, cmdIdle); err != nil {
		return 0, err
	}

	valH, err := dev.ReadRegister(hw.CRCResultRegH)
	if err != nil {
		return 0, err
	}
	valL, err := dev.ReadRegister(hw.CRCResultRegL)
	if err != nil {
		return 0, err
	}
	return uint16(valH) << 8 | uint16(valL), nil
}
*/

func SoftReset(dev DeviceIO) error {
	if err := command(dev, cmdSoftReset); err != nil {
		return err
	}
	time.Sleep(50 * time.Millisecond)
	return waitRegister(dev, hw.CommandReg, 0x00, hw.PowerDown, 1 * time.Second)
}
