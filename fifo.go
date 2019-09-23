package rc522

import (
	"fmt"
	"time"
	"errors"
	"github.com/op0xA5/rc522/hw"
)

func flushFIFO(dev DeviceIO) error {
	err := dev.WriteRegister(hw.FIFOLevelReg, hw.FlushBuffer)
	if err != nil {
		return fmt.Errorf("flush fifo: %v", err)
	}
	return nil
}

func readFromFIFO(dev DeviceIO, n int, timeout time.Duration) ([]byte, error) {
	buffer := make([]byte, n)
	pos := 0

	deadline := newDeadline(timeout)

	for {
		n, err := dev.ReadFIFO(buffer[pos:])
		fmt.Println(n, buffer)
		pos = pos + n
		if err != nil {
			return buffer[:pos], err
		}
		if pos == len(buffer) {
			return buffer, nil
		}
		if deadline.Timeout() {
			return buffer[:pos], errors.New("timeout")
		}
		time.Sleep(50 * time.Millisecond)
	}
}

func writeToFIFO(dev DeviceIO, data []byte, timeout time.Duration) error {
	deadline := newDeadline(timeout)

	for {
		n, err := dev.WriteFIFO(data)
		if err != nil {
			return err
		}
		if n == len(data) {
			return nil
		}
		if deadline.Timeout() {
			return errors.New("timeout")
		}
		data = data[n:]
		time.Sleep(50 * time.Millisecond)
	}
}