package rc522

import (
	"time"
	"errors"
	"fmt"
)

type BlockTransfer struct {
	CID        uint8
	FrameSize  int
	FWT        time.Duration

	SendBlockIndex    int
	ReceiveBlockIndex int
}

const (
	pcbBlockI = 0x02
	pcbBlockR_ACK = 0xA2
	pcbBlockR_NAK = 0xB2
	pcbDeslect = 0xCA
	pcbWithCID = 0x08
	pcbLink    = 0x10
)

func (bt *BlockTransfer) Send(dev DeviceIO, data []byte) ([]byte, error) {
	var buffer = make([]byte, bt.FrameSize)
	var recvBuffer = make([]byte, 0, 2 * 1024)

	frameDataSize := bt.FrameSize - 1 - 1 - 2

	var pcb byte
	var recvData []byte

	// send loop
	for {
		req := TransceiveRequest()

		isLinkReq := len(data) > frameDataSize
		buffer[0] = pcbBlockI | pcbWithCID | byte(bt.SendBlockIndex & 0x01)
		if isLinkReq {
			buffer[0] |= pcbLink
		}
		buffer[1] = bt.CID
		sendSize := copy(buffer[2:2+frameDataSize], data)
		req.SendData = buffer[:2+sendSize]
		req.SendData = AppendISO14443aCRC(req.SendData, req.SendData)
		req.CheckCRC = true
		fmt.Printf("block: send %v\n", req.SendData)
		res, err := req.Send(dev)
		if err != nil {
			return nil, err
		}
		fmt.Printf("block: recv %v\n", res.ReceiveData)
		if len(res.ReceiveData) < 1 {
			return nil, errors.New("bad picc response")
		}
		pcb = res.ReceiveData[0]

		if !isLinkReq {
			recvData = res.ReceiveData
			bt.SendBlockIndex++
			break
		}

		if pcb & 0xC0 != 0x80 /* R-block */ {
			return nil, errors.New("bad picc response")
		}
		if pcb & 0x10 != 0 /* ACK */{
			return nil, errors.New("picc NAK")
		}

		if pcb & 0x01 != uint8(bt.SendBlockIndex & 0x01) {
			// retry send block
			continue
		}

		bt.SendBlockIndex++
		data = data[sendSize:]
	}

	// receive loop
	for {
		if pcb & 0xC0 != 0x00 /* I-block */ {
			return recvBuffer, errors.New("bad picc response")
		}

		if pcb & 0x01 == uint8(bt.ReceiveBlockIndex & 0x01) {
			payload := recvData[1:]
			if pcb & pcbWithCID != 0 {
				payload = payload[1:]
			}
			recvBuffer = append(recvBuffer, payload...)
			bt.ReceiveBlockIndex++
		}

		if pcb & pcbLink == 0 {
			// not link
			break
		}

		buffer[0] = pcbBlockR_ACK | pcbWithCID | byte(bt.ReceiveBlockIndex & 0x01)
		buffer[1] = bt.CID
		AppendISO14443aCRC(buffer[2:2], buffer[:2])
		req := TransceiveRequest()
		req.SendData = buffer[:4]
		req.CheckCRC = true
		fmt.Printf("block: send %v\n", req.SendData)
		res, err := req.Send(dev)
		if err != nil {
			return recvBuffer, err
		}
		fmt.Printf("block: recv %v\n", res.ReceiveData)
		if len(res.ReceiveData) < 1 {
			return nil, errors.New("bad picc response")
		}
		pcb = res.ReceiveData[0]
		recvData = res.ReceiveData
	}

	return recvBuffer, nil
}
