package rc522

import (
	"testing"
)

func TestISO14443aCRC(t *testing.T) {
	data := []byte{0x50, 0x00}
	result := ISO14443aCRC(data)
	t.Logf("CRC Result: %04x", result)

	if result != 0xCD57 { // shoule be "57 cd" in packet
		t.Fatal("CRC result shoule be 0xCD57")
	}

	buffer := []byte{0x50, 0x00, 0x00, 0x00}
	AppendISO14443aCRC(buffer[2:2], buffer[:2])
	t.Logf("Append CRC: %v", buffer)
}
