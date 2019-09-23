package hw

const (
	CmdIdle             = uint8(0x00)
	CmdMem              = uint8(0x01)
	CmdGenerateRandomID = uint8(0x02)
	CmdCalcCRC          = uint8(0x03)
	CmdTransmit         = uint8(0x04)
	CmdNoCmdChange      = uint8(0x07)
	CmdReceive          = uint8(0x08)
	CmdTransceive       = uint8(0x0C)
	CmdMFAuthent        = uint8(0x0E)
	CmdSoftReset        = uint8(0x0F)
)

func CommandNeedsRX(cmd uint8) (needsRx bool, ok bool) {
	switch (cmd) {
	case CmdIdle, CmdMem, CmdGenerateRandomID, CmdCalcCRC,
		CmdTransmit, CmdSoftReset:
		return false, true
	case CmdReceive, CmdTransceive, CmdMFAuthent:
		return true, true
	}
	return false, false
}