package rc522

import (
	"fmt"

	rpio "github.com/stianeikeland/go-rpio"
	"github.com/op0xA5/rc522/hw"
)

// DeviceIO SI522/RC522 IO通信模型
type DeviceIO interface {
	ReadRegister(addr uint8) (uint8, error)
	WriteRegister(addr uint8, value uint8) error
	ReadFIFO(p []byte) (n int, err error)
	WriteFIFO(p []byte) (n int, err error)
	Close() error
}

type deviceViaSPI struct {
	spiDev rpio.SpiDev
	csPin  rpio.Pin
}

const (
	csPinHardware = 0xff
)
const (
	SpiCE0 = rpio.Pin(8)
	SpiCE1 = rpio.Pin(7)
)

func OpenSPI(spiDev rpio.SpiDev, csPin rpio.Pin, speed int) (DeviceIO, error) {
	if err := rpio.Open(); err != nil {
		return nil, fmt.Errorf("ERR: rpio.Open, %v", err)
	}

	if err := rpio.SpiBegin(spiDev); err != nil {
		return nil, fmt.Errorf("ERR: rpio.SpiBegin, %v", err)
	}

	switch csPin {
	case SpiCE0:
		rpio.SpiChipSelect(0)
		rpio.SpiChipSelectPolarity(0, 0)
		csPin = csPinHardware
	case SpiCE1:
		rpio.SpiChipSelect(1)
		rpio.SpiChipSelectPolarity(1, 0)
		csPin = csPinHardware
	}

	if csPin != csPinHardware {
		csPin.High()
		csPin.Output()
		csPin.High()
	}

	rpio.SpiSpeed(speed)
	rpio.SpiMode(0, 0)

	return &deviceViaSPI{spiDev, csPin}, nil
}
const (
	addressRead = uint8(0x80)
)
func (dev *deviceViaSPI) ReadRegister(addr uint8) (uint8, error) {
	var buffer [2]byte
	buffer[0] = (addr << 1) | addressRead
	buffer[1] = 0x00
	if dev.csPin == csPinHardware {
		rpio.SpiExchange(buffer[:])
	} else {
		dev.csPin.Low()
		rpio.SpiExchange(buffer[:])
		dev.csPin.High()
	}
	return buffer[1], nil
}
func (dev *deviceViaSPI) WriteRegister(addr uint8, value uint8) error {
	var buffer [2]byte
	buffer[0] = (addr << 1) & 0x7F
	buffer[1] = value
	if dev.csPin == csPinHardware {
		rpio.SpiExchange(buffer[:])
	} else {
		dev.csPin.Low()
		rpio.SpiExchange(buffer[:])
		dev.csPin.High()
	}
	return nil
}
func (dev *deviceViaSPI) ReadFIFO(p []byte) (n int, err error) {
	fifoLevel, err := dev.ReadRegister(hw.FIFOLevelReg)
	if err != nil {
		return 0, err
	}

	dataSize := int(fifoLevel & hw.FIFOLevel)
	if dataSize > hw.FIFOMaxSize {
		dataSize = hw.FIFOMaxSize
	}
	if dataSize > len(p) {
		dataSize = len(p)
	}

	var buffer [hw.FIFOMaxSize+1]byte
	for i := 0; i < dataSize; i++ {
		buffer[i] = (hw.FIFODataReg << 1) | addressRead
	}
	buffer[dataSize] = 0x00

	if dev.csPin == csPinHardware {
		rpio.SpiExchange(buffer[:dataSize+1])
	} else {
		dev.csPin.Low()
		rpio.SpiExchange(buffer[:dataSize+1])
		dev.csPin.High()
	}

	copy(p, buffer[1:dataSize+1])

	return dataSize, nil
}
func (dev *deviceViaSPI) WriteFIFO(p []byte) (n int, err error) {
	fifoLevel, err := dev.ReadRegister(hw.FIFOLevelReg)
	if err != nil {
		return 0, err
	}

	dataSize := hw.FIFOMaxSize - int(fifoLevel & hw.FIFOLevel)
	if dataSize <= 0 {
		return 0, nil
	}
	if dataSize > len(p) {
		dataSize = len(p)
	}

	var buffer [hw.FIFOMaxSize+1]byte
	buffer[0] = (hw.FIFODataReg << 1)
	copy(buffer[1:dataSize+1], p)

	if dev.csPin == csPinHardware {
		rpio.SpiExchange(buffer[:dataSize+1])
	} else {
		dev.csPin.Low()
		rpio.SpiExchange(buffer[:dataSize+1])
		dev.csPin.High()
	}

	return dataSize, nil
}
func (dev *deviceViaSPI) Close() error {
	rpio.SpiEnd(dev.spiDev)
	if dev.csPin != csPinHardware {
		dev.csPin.Input()
	}
	return rpio.Close()
}

