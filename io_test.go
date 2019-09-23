package rc522

import (
	"time"
	"testing"
)

func TestOpenSPI(t *testing.T) {
	dev, err := OpenSPI(0, SpiCE0, 5 * 1000000)
	if err != nil {
		t.Fatal(err)
		return
	}

	t.Log("device opened")

	err = dev.Close()
	if err != nil {
		t.Fatal(err)
	}

	t.Log("device closed")
}

func TestReadRegister(t *testing.T) {
	dev, err := OpenSPI(0, SpiCE0, 5 * 1000000)
	if err != nil {
		t.Fatal(err)
		return
	}

	t.Log("device opened")

	v, err := dev.ReadRegister(0x37 /*37h: VersionReg*/)
	if err != nil {
		t.Fatal(err)
		return
	}

	switch v {
	case 0x91:
		t.Log("chip is MFRC522 version 1.0")
	case 0x92:
		t.Log("chip is MFRC522 version 2.0")
	default:
		t.Logf("unknown chip id: %x", v)
	}

	err = dev.Close()
	if err != nil {
		t.Fatal(err)
	}

	t.Log("device closed")
}

func TestReadFIFO(t *testing.T) {
	dev, err := OpenSPI(0, SpiCE0, 5 * 1000000)
	if err != nil {
		t.Fatal(err)
		return
	}

	t.Log("device opened")

	err = flushFIFO(dev)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("write 16 bytes to FIFO")
	data := []byte{ 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15 }
	err = writeToFIFO(dev, data, 1 * time.Second)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("read 25 bytes from FIFO")
	p, err := readFromFIFO(dev, 25, 1 * time.Second)

	t.Logf("%d bytes read: %v", len(p), p)

	err = dev.Close()
	if err != nil {
		t.Fatal(err)
	}

	t.Log("device closed")
}

func TestGenerateRandom(t *testing.T) {
	dev, err := OpenSPI(0, SpiCE0, 5 * 1000000)
	if err != nil {
		t.Fatal(err)
		return
	}

	t.Log("device opened")

	err = flushFIFO(dev)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("cmdGenerateRandomID")
	err = command(dev, cmdGenerateRandomID)
	if err != nil {
		t.Fatal(err)
	}

	err = waitIdle(dev,1 * time.Second)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("cmdMem")
	err = command(dev, cmdMem)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("read 25 bytes from FIFO")
	p, err := readFromFIFO(dev, 25, 1 * time.Second)

	t.Logf("%d bytes read: %v", len(p), p)

	err = dev.Close()
	if err != nil {
		t.Fatal(err)
	}

	t.Log("device closed")
}
