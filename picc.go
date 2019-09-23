package rc522

import (
	"time"
	"errors"
	"fmt"
	"github.com/op0xA5/rc522/hw"
)

const (
	PiccCmdREQA   = uint8(0x26)
	PiccCmdWUPA   = uint8(0x52)
	PiccCmdCT     = uint8(0x88)
	PiccCmdSELCL1 = uint8(0x93)
	PiccCmdSELCL2 = uint8(0x95)
	PiccCmdSELCL3 = uint8(0x97)
	PiccCmdRATS   = uint8(0xE0)
	PiccCmdHLTA   = uint8(0x50)
)

const (
	PiccTypeUnknown    = iota
	PiccTypeISO14443_4
	PiccTypeISO18092
	PiccTypeMifareMini
	PiccTypeMifare1K
	PiccTypeMifare4K
	PiccTypeMifareUL
	PiccTypeMifarePlus
	PiccTypeNotComplete = 255
)

type ATQA uint16
type SAK  uint8
type UID  []byte

func RequestA(dev DeviceIO) (ATQA, error) {
	return reqa_or_wupa(dev, PiccCmdREQA)
}
func WakeupA(dev DeviceIO) (ATQA, error) {
	return reqa_or_wupa(dev, PiccCmdWUPA)
}
func reqa_or_wupa(dev DeviceIO, piccCmd uint8) (ATQA, error) {
	req := TransceiveRequest()
	var sendData [1]byte
	sendData[0] = piccCmd
	req.SendData = sendData[:]
	req.ValidBits = 7
	req.Timeout = 250 * time.Millisecond
	resp, err := req.Send(dev)
	if err != nil {
		return 0, err
	}
	if len(resp.ReceiveData) != 2 || resp.ValidBits != 0 {
		return 0, errors.New("invalid response")
	}

	return (ATQA(resp.ReceiveData[0]) << 8) | ATQA(resp.ReceiveData[1]), nil
}

func (atqa ATQA) ProprietaryCoding() int {
	return int(atqa >> 8) & 0x0F
}
func (atqa ATQA) UIDSize() int {
	return int(atqa >> 6) & 0x03
}
func (atqa ATQA) UIDBytes() (int, error) {
	switch atqa.UIDSize() {
	case 0:
		return 4, nil
	case 1:
		return 7, nil
	case 2:
		return 10, nil
	}
	return 0, errors.New("unknown UID bytes")
}
func (atqa ATQA) BitFrameAnticollision() int {
	return int(atqa) & 0x1F
}

func Select(dev DeviceIO) (UID, SAK, error) {
	var uidBuffer [10]byte
	var uid []byte
	var sak SAK

	req := TransceiveRequest()
	var buffer [9]byte

	for cascadeLevel := 1; cascadeLevel <= 3; cascadeLevel++ {
		var currentLevelUid []byte

		switch cascadeLevel {
		case 1:
			buffer[0] = PiccCmdSELCL1
			currentLevelUid = uidBuffer[:]
		case 2:
			buffer[0] = PiccCmdSELCL2
			currentLevelUid = uidBuffer[3:]
		case 3:
			buffer[0] = PiccCmdSELCL3
			currentLevelUid = uidBuffer[6:]
		}

		buffer[1] = 0x20

		req.SendData = buffer[:2]
		req.CheckCRC = false
		res, err := req.Send(dev)
		if err != nil {
			return nil, 0, err
		}
		if len(res.ReceiveData) < 5 {
			return nil, 0, errors.New("bad picc response")
		}

		if res.ReceiveData[0] == PiccCmdCT {
			copy(currentLevelUid, res.ReceiveData[1:])
			uid = uidBuffer[:len(uid) + 3]
		} else {
			copy(currentLevelUid, res.ReceiveData)
			uid = uidBuffer[:len(uid) + 4]
		}

		// Select command
		buffer[1] = 0x70
		copy(buffer[2:7], res.ReceiveData)
		AppendISO14443aCRC(buffer[7:7], buffer[:7])
		req.SendData = buffer[:9]
		req.CheckCRC = true
		res, err = req.Send(dev)
		if err != nil {
			return nil, 0, err
		}
		if len(res.ReceiveData) != 1 {
			return nil, 0, errors.New("bad picc response")
		}
		sak = SAK(res.ReceiveData[0])

		if sak.UIDComplete() {
			return uid, sak, nil
		}
	}

	return nil, 0, errors.New("bad SAK response")
}

// TODO: finish anticollision logic
func anticollision(dev DeviceIO, cascadeCmd uint8, p []byte, knownBits int) error {
	panic("not implements")

	if len(p) < 5 {
		return errors.New("anticollision input/output buffer size shoule be 5 bytes")
	}
	if knownBits < 0 {
		return errors.New("bad knownBits value")
	}
	if knownBits >= 4 * 8 {
		return errors.New("bad knownBits value, may should use select command")
	}

	req := TransceiveRequest()
	req.ValuesAfterColl = true

	var buffer [9]byte
	buffer[0] = cascadeCmd

	for {
		buffer[1] = 0x20 + uint8((knownBits / 8) << 4) + uint8(knownBits % 8)
		copy(buffer[2:], p)
		req.SendData = buffer[:2+(knownBits+7)/8]
		req.ValidBits = knownBits % 8
		req.RxAlign = knownBits % 8
		res, err := req.Send(dev)

		if res != nil && len(res.ReceiveData) > 0 {			
			mask := uint8(0xFF << knownBits % 8) // new receive part
			bytePos := knownBits / 8
			res.ReceiveData[0] &= mask
			p[bytePos] &= ^mask
			p[bytePos] |= res.ReceiveData[0]
			copy(p[bytePos+1:], res.ReceiveData[1:])
		}

		if err != nil {
			if IsCollision(err) {
				v, err := dev.ReadRegister(hw.CollReg)
				if err != nil {
					return err
				}
				collpos := int(v & hw.CollPos)
				if collpos == 0 {
					collpos = 32
				}
				fmt.Printf("coll pos: %v\n", collpos)
				knownBits += collpos - 1
				p[knownBits / 8] |= 1 << (knownBits % 8)
				knownBits++
				continue
			}
			return err
		}
		break
	}

	return nil
}

func (sak SAK) UIDComplete() bool {
	return sak & 0x04 == 0
}
func (sak SAK) Support14443_4() bool {
	return sak & 0x24 == 0x20
}
func (sak SAK) NotSupport14443_4() bool {
	return sak & 0x24 == 0x00
}

type FrameSize int
const FrameSizeRFU = -1
var frameSizeTbl = [16]FrameSize{
	16, 24, 32, 40, 48, 64, 96, 128, 256,
	FrameSizeRFU, FrameSizeRFU, FrameSizeRFU, FrameSizeRFU,
	FrameSizeRFU, FrameSizeRFU, FrameSizeRFU,
}
func frameSize(fsi uint8) FrameSize {
	return frameSizeTbl[fsi & 0x0F]
}
func (fsc FrameSize) Index() uint8 {
	switch {
	case fsc <= 16:
		return 0
	case fsc <= 24:
		return 1
	case fsc <= 32:
		return 2
	case fsc <= 40:
		return 3
	case fsc <= 48:
		return 4
	case fsc <= 64:
		return 5
	case fsc <= 96:
		return 6
	case fsc <= 128:
		return 7
	default:
		return 8
	}
}

type ATS []byte
func (ats ATS) FrameSize() FrameSize {
	return frameSize(ats[1])
}

func RequestATS(dev DeviceIO, fs FrameSize, cid uint8) (ATS, error) {
	var buffer [4]byte
	buffer[0] = PiccCmdRATS
	buffer[1] = (fs.Index() << 4) | (cid & 0x0F)
	AppendISO14443aCRC(buffer[2:2], buffer[:2])
	req := TransceiveRequest()
	req.SendData = buffer[:]
	req.CheckCRC = true
	res, err := req.Send(dev)
	if err != nil {
		return nil, err
	}
	return ATS(res.ReceiveData), nil
}

func HaltA(dev DeviceIO) error {
	req := TransceiveRequest()
	var sendData [4]byte
	sendData[0] = PiccCmdHLTA
	sendData[1] = 0
	sendData[2] = 0x57  // CRC
	sendData[3] = 0xCD
	req.SendData = sendData[:]
	req.Timeout = 250 * time.Millisecond
	_, err := req.Send(dev)
	return err
}
