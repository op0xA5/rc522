package rc522

import (
	"time"
	"testing"
)

func TestRequestA(t *testing.T) {
	dev, err := OpenSPI(0, SpiCE0, 1 * 1000000)
	if err != nil {
		t.Fatal(err)
		return
	}

	t.Log("device opened")

	err = InitPCD(dev)
	if err != nil {
		t.Fatal(err)
		return
	}

	time.Sleep(1 * time.Millisecond)

	t.Logf("send REQA")

	atqa, err := RequestA(dev)
	t.Logf("err: %v", err)
	t.Logf("ATQA: %04x", atqa)
	
	uid, sak, err := Select(dev)
	t.Logf("err: %v", err)
	t.Logf("uid: %v", uid)
	t.Logf("sak: %02x", sak)

	ats, err := RequestATS(dev, 64, 1)
	t.Logf("err: %v", err)
	t.Logf("ats: %v", ats)

	t.Logf("FSC: %v", ats.FrameSize())

	bt := new(BlockTransfer)
	bt.CID = 1
	bt.FrameSize = 64

	resp, err := bt.Send(dev, []byte{ 0x00, 0x84, 0x00, 0x00, 8 })
	t.Logf("err: %v", err)
	t.Logf("GET CHALLENGE resp: %v", resp)

	resp, err = bt.Send(dev, []byte{ 0x00, 0x84, 0x00, 0x00, 8 })
	t.Logf("err: %v", err)
	t.Logf("GET CHALLENGE resp: %v", resp)

	resp, err = bt.Send(dev, []byte{ 0x00, 0x84, 0x00, 0x00, 8 })
	t.Logf("err: %v", err)
	t.Logf("GET CHALLENGE resp: %v", resp)

	AntennaOff(dev)

	err = dev.Close()
	if err != nil {
		t.Fatal(err)
	}

	t.Log("device closed")
}
