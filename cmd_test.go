package rc522

import (
	"testing"
)

func TestSoftReset(t *testing.T) {
	dev, err := OpenSPI(0, SpiCE0, 5 * 1000000)
	if err != nil {
		t.Fatal(err)
		return
	}

	t.Log("device opened")

	err = SoftReset(dev)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("soft reset performed")

	err = dev.Close()
	if err != nil {
		t.Fatal(err)
	}

	t.Log("device closed")
}
