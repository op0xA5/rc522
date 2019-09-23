package hw

const (
	CommandReg    = 0x01
	ComlEnReg     = 0x02
	DivlEnReg     = 0x03
	ComIrqReg     = 0x04
	DivIrqReg     = 0x05
	ErrorReg      = 0x06
	Status1Reg    = 0x07
	Status2Reg    = 0x08
	FIFODataReg   = 0x09
	FIFOLevelReg  = 0x0A
	WaterLevelReg = 0x0B
	ControlReg    = 0x0C
	BitFramingReg = 0x0D
	CollReg       = 0x0E

	ModeReg        = 0x11
	TxModeReg      = 0x12
	RxModeReg      = 0x13
	TxControlReg   = 0x14
	TxASKReg       = 0x15
	TxSelReg       = 0x16
	RxSelReg       = 0x17
	RxThresholdReg = 0x18
	DemodReg       = 0x19
	MfTxReg        = 0x1C
	MfRxReg        = 0x1D
	SerialSpeedReg = 0x1F

	CRCResultRegH   = 0x21
	CRCResultRegL   = 0x22
	ModWidthReg     = 0x24
	RFCfgReg        = 0x26
	GsNReg          = 0x27
	CWGsPReg        = 0x28
	ModGsPReg       = 0x29
	TModeReg        = 0x2A
	TPrescalerReg   = 0x2B
	TReloadRegH     = 0x2C
	TReloadRegL     = 0x2D
	TCounterValRegH = 0x2E
	TCounterValRegL = 0x2F

	TestSel1Reg     = 0x31
	TestSel2Reg     = 0x32
	TestPinEnReg    = 0x33
	TestPinValueReg = 0x34
	TestBusReg      = 0x35
	AutoTestReg     = 0x36
	VersionReg      = 0x37
	AnalogTestReg   = 0x38
	TestDAC1Reg     = 0x39
	TestDAC2Reg     = 0x3A
	TestADCReg      = 0x3B
)


const (
	RcvOff    = uint8(1 << 5)
	PowerDown = uint8(1 << 4)
	Command   = uint8(0x0F)

	IRqInv     = uint8(1 << 7)
	TxIEn      = uint8(1 << 6)
	RxIEn      = uint8(1 << 5)
	IdleIEn    = uint8(1 << 4)
	HiAlertIEn = uint8(1 << 3)
	LoAlertIEn = uint8(1 << 2)
	ErrIEn     = uint8(1 << 1)
	TimerIEn   = uint8(1 << 0)

	IRQPushPull = uint8(1 << 7)
	MfinActIEn  = uint8(1 << 4)
	CRCIEn      = uint8(1 << 2)

	Set1       = uint8(1 << 7)
	TxIRq      = uint8(1 << 6)
	RxIRq      = uint8(1 << 5)
	IdleIRq    = uint8(1 << 4)
	HiAlertIRq = uint8(1 << 3)
	LoAlertIRq = uint8(1 << 2)
	ErrIRq     = uint8(1 << 1)
	TimerIRq   = uint8(1 << 0)

	Set2       = uint8(1 << 7)
	MfinActIRq = uint8(1 << 4)
	CRCIRq     = uint8(1 << 2)

	WrErr       = uint8(1 << 7)
	TempErr     = uint8(1 << 6)
	BufferOvfl  = uint8(1 << 4)
	CollErr     = uint8(1 << 3)
	CRCErr      = uint8(1 << 2)
	ParityErr   = uint8(1 << 1)
	ProtocolErr = uint8(1 << 0)

	CRCOk    = uint8(1 << 6)
	CRCReady = uint8(1 << 5)
	IRq      = uint8(1 << 4)
	TRunning = uint8(1 << 3)
	HiAlert  = uint8(1 << 1)
	LoAlert  = uint8(1 << 0)

	TempSensClear = uint8(1 << 7)
	I2CForceHS    = uint8(1 << 6)
	MFCrypto1On   = uint8(1 << 3)
	ModenState    = uint8(0x07)
	ModemState2   = uint8(1 << 2)
	ModemState1   = uint8(1 << 1)
	ModemState0   = uint8(1 << 0)

	FlushBuffer = uint8(1 << 7)
	FIFOLevel   = uint8(0x7F)

	WaterLevel = uint8(0x3F)

	TStopNow   = uint8(1 << 7)
	TStartNow  = uint8(1 << 6)
	RxLastBits = uint8(0x07)
 
	StartSend  = uint8(1 << 7)
	RxAlign    = uint8(0x70)
	TxLastBits = uint8(0x07)

	ValuesAfterColl = uint8(1 << 7)
	CollPosNotValid = uint8(1 << 5)
	CollPos         = uint8(0x1F)

	MSBFirst  = uint8(1 << 7)
	TxWaitRF  = uint8(1 << 5)
	PolMFin   = uint8(1 << 3)
	CRCPreset = uint8(0x03)

	TxCRCEn = uint8(1 << 7)
	TxSpeed = uint8(0x70)
	InvMod  = uint8(1 << 3)

	RxCRCEn    = uint8(1 << 7)
	RxSpeed    = uint8(0x70)
	RxNoErr    = uint8(1 << 3)
	RxMultiple = uint8(1 << 2)

	InvTx2RFOn  = uint8(1 << 7)
	InvTx1RFOn  = uint8(1 << 6)	
	InvTx2RFOff = uint8(1 << 5)
	InvTx1RFOff = uint8(1 << 4)
	Tx2CW       = uint8(1 << 3)
	Tx2RFEn     = uint8(1 << 1)
	Tx1RFEn     = uint8(1 << 0)

	Force100ASK = uint8(1 << 6)

	DriverSel = uint8(0x30)
	MFOutSel  = uint8(0x0F)

	UARTSel = uint8(0xC0)
	RxWait  = uint8(0x3F)

	MinLevel  = uint8(0xF0)
	CollLevel = uint8(0x07)

	AddIQ        = uint8(0xC0)
	FixIQ        = uint8(1 << 5)
	TPrescalEven = uint8(1 << 4)
	TauRcv       = uint8(0x0C)
	TauSync      = uint8(0x03)

	TxWait = uint8(0x03)

	ParityDisable = uint8(1 << 4)

	BR_T0 = uint8(0xE0)
	BR_T1 = uint8(0x1F)

	RxGain = uint8(0x70)

	CWGsN  = uint8(0xF0)
	ModGsN = uint8(0x0F)
	CWGsP  = uint8(0x3F)
	ModGsP = uint8(0x3F)

	TAuto        = uint8(1 << 7)
	TGated       = uint8(0x60)
	TAutoRestart = uint8(1 << 4)
	TPrescalerH  = uint8(0x0F)
	TPrescalerL  = uint8(0xFF)
)

const (
	TxSpeed106kBd = 0x00
	TxSpeed212kBd = 0x01
	TxSpeed424kBd = 0x02
	TxSpeed848kBd = 0x03
)
