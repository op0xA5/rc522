package rc522

import (
	"time"
	"errors"
	"fmt"
	"github.com/op0xA5/rc522/hw"
)

func InitPCD(dev DeviceIO) error {
	if err := SoftReset(dev); err != nil {
		return err
	}

	if err := dev.WriteRegister(hw.TxASKReg, hw.Force100ASK); err != nil {
		return err
	}
	// Default 0x3F. Set the preset value for the CRC coprocessor for the CalcCRC command to 0x6363 (ISO 14443-3 part 6.2.4)
	if err := dev.WriteRegister(hw.ModeReg, 0x3D); err != nil {
		return err
	}
/*
	if err := dev.WriteRegister(hw.TModeReg, 0x80); err != nil {
		return err
	}
	if err := dev.WriteRegister(hw.TPrescalerReg, 0xA9); err != nil {
		return err
	}
	if err := dev.WriteRegister(hw.TReloadRegH, 0x03); err != nil {
		return err
	}
	if err := dev.WriteRegister(hw.TReloadRegL, 0xE8); err != nil {
		return err
	}

	if err := dev.WriteRegister(hw.RxSelReg, 0x86); err != nil {
		return err
	}
*/
	return AntennaOn(dev)
}
func AntennaOn(dev DeviceIO) error {
	val, err := dev.ReadRegister(hw.TxControlReg)
	if err != nil {
		return err
	}
	if val & (hw.Tx1RFEn | hw.Tx2RFEn) == (hw.Tx1RFEn | hw.Tx2RFEn) {
		return nil
	}
	val |= hw.Tx1RFEn | hw.Tx2RFEn
	return dev.WriteRegister(hw.TxControlReg, val)
}
func AntennaOff(dev DeviceIO) error {
	val, err := dev.ReadRegister(hw.TxControlReg)
	if err != nil {
		return err
	}
	val &= ^(hw.Tx1RFEn | hw.Tx2RFEn)
	return dev.WriteRegister(hw.TxControlReg, val)
}

type AntennaGain int
const (
	AntennaGain18dB = AntennaGain(18)
	AntennaGain23dB = AntennaGain(23)
	AntennaGain33dB = AntennaGain(33)
	AntennaGain38dB = AntennaGain(38)
	AntennaGain43dB = AntennaGain(43)
	AntennaGain48dB = AntennaGain(48)
)
func GetAntennaGain(dev DeviceIO) (AntennaGain, error) {
	val, err := dev.ReadRegister(hw.RFCfgReg)
	if err != nil {
		return AntennaGain(0), err
	}
	rxGain := (val & hw.RxGain) >> 4
	switch rxGain {
	case 0:
		return AntennaGain(18), err
	case 1:
		return AntennaGain(23), err
	case 2:
		return AntennaGain(18), err
	case 3:
		return AntennaGain(23), err
	case 4:
		return AntennaGain(33), err
	case 5:
		return AntennaGain(38), err
	case 6:
		return AntennaGain(43), err
	case 7:
		return AntennaGain(48), err
	}
	return AntennaGain(33), nil
}
func SetAntennaGain(dev DeviceIO, ag AntennaGain) error {
	var rxGain = uint8(0x40)
	switch {
	case ag <= 18:
		rxGain = 0x00
	case ag <= 23:
		rxGain = 0x10
	case ag <= 33:
		rxGain = 0x40
	case ag <= 38:
		rxGain = 0x50
	case ag <= 43:
		rxGain = 0x60
	default:
		rxGain = 0x70
	}
	val, err := dev.ReadRegister(hw.RFCfgReg)
	if err != nil {
		return err
	}
	if val & hw.RxGain != rxGain {
		val &= ^hw.RxGain
		val |= rxGain
		return dev.WriteRegister(hw.RFCfgReg, val)
	}
	return nil
}

type Request struct {
	// Command pass to PCD chip (RC522)
	// if dont know which command to send, use TransceiveConversation of other func to setup a Conversation
	Command         int
	SendData        []byte
	ValidBits       int
	RxAlign         int
	CheckCRC        bool
	Timeout         time.Duration
	ValuesAfterColl bool
}

type Response struct {
	ValidBits   int
	ReceiveData []byte
}

func printComIrq(dev DeviceIO) {
	val, err := dev.ReadRegister(hw.ComIrqReg)
	fmt.Printf("printComIrq: %02x  %v\n", val, err)
}
func printFIFOLevel(dev DeviceIO) {
	val, err := dev.ReadRegister(hw.FIFOLevelReg)
	fmt.Printf("printFIFOLevel: %02x  %v\n", val, err)
}
func printControlReg(dev DeviceIO) {
	val, err := dev.ReadRegister(hw.ControlReg)
	fmt.Printf("printControlReg: %02x  %v\n", val, err)
}

func TransceiveRequest() *Request {
	return &Request{
		Command: int(cmdTransceive),
		Timeout: 250 * time.Millisecond,
	}
}

func (r *Request) Send(dev DeviceIO) (*Response, error) {
	deadline := newDeadline(r.Timeout)

	recvCmd := r.Command == int(cmdTransceive) || r.Command == int(cmdReceive)

	if len(r.SendData) > 64 {
		return nil, errors.New("send data large than 64bytes not support currently")
	}

	if err := command(dev, cmdIdle); err != nil {
		return nil, err
	}
	if err := irqCom.clear(dev); err != nil {
		return nil, err
	}

	if recvCmd {
		val, err := dev.ReadRegister(hw.CollReg)
		if err != nil {
			return nil, err
		}
		nval := val
		if r.ValuesAfterColl {
			nval = nval | hw.ValuesAfterColl
		} else {
			nval = nval & ^hw.ValuesAfterColl
		}
		if val != nval {
			err := dev.WriteRegister(hw.CollReg, nval)
			if err != nil {
				return nil, err
			}
		}
	}

	if err := flushFIFO(dev); err != nil {
		return nil, err
	}
	if len(r.SendData) > 0 {
		n, err := dev.WriteFIFO(r.SendData)
		if err != nil {
			return nil, err
		}
		if n != len(r.SendData) {
			return nil, errors.New("short write")
		}
	}

	var bitFraming = uint8((r.RxAlign & 0x07) << 4 | r.ValidBits & 0x07)
	if err := dev.WriteRegister(hw.BitFramingReg, bitFraming); err != nil {
		return nil, err
	}
	if err := command(dev, uint8(r.Command)); err != nil {
		return nil, err
	}
	
	if r.Command == int(cmdTransceive) || r.Command == int(cmdTransmit) {
		bitFraming |= hw.StartSend
		if err := dev.WriteRegister(hw.BitFramingReg, bitFraming); err != nil {
			return nil, err
		}
	}
	
	waitIrq := irqIdle | irqErr
	if recvCmd {
		waitIrq |= irqRx
	}
	if err := waitIrq.wait(dev, deadline.SinceNow()); err != nil {
		return nil, err
	}
	
	var errorReg uint8
	var err error
	if errorReg, err = dev.ReadRegister(hw.ErrorReg); err != nil {
		return nil, err
	}

	if errorReg & (hw.BufferOvfl | hw.ParityErr | hw.ProtocolErr) > 0 {
		return nil, errors.New("receive error")
	}

	var response Response
	var recvCrc  uint16
	if recvCmd {
		if err := irqRx.wait(dev, deadline.SinceNow()); err != nil {
			return nil, err
		}

		var p [64]byte
		n, err := dev.ReadFIFO(p[:])
		if err != nil {
			return &response, err
		}

		var val uint8
		if val, err = dev.ReadRegister(hw.ControlReg); err != nil {
			return &response, err
		}
		response.ValidBits = int(val & hw.RxLastBits)

		if r.CheckCRC {
			if n < 2 || response.ValidBits != 0 {
				response.ReceiveData = nil
				recvCrc = 0
			} else {
				response.ReceiveData = p[:n-2]
				recvCrc = (uint16(p[n-2]) << 8 | uint16(p[n-1]))
			}
		} else {
			response.ReceiveData = p[:n]
		}
	}

	if errorReg & hw.CollErr > 0 {
		return &response, collisionErr
	}

	if recvCmd && r.CheckCRC {
		if response.ReceiveData == nil && recvCrc == 0 {
			return &response, badCrcErr
		}

		crc := ISO14443aCRC(response.ReceiveData)
		if recvCrc != crc {
			fmt.Printf("crc: %04x, got: %04x\n", crc, recvCrc)
			return &response, badCrcErr
		}
	}

	return &response, nil
}

