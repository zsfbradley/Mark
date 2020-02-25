package mFace

type MDataProtocolStruct interface {
	Symbol() string
	LengthOfStages() int
	ValueOfStages(uint8) uint64
	AppendData([]byte)
	CleanData()
	Data() [][]byte
	SymbolLength() uint64
	HeaderLength() uint64
	RangeOfDataLength() (uint64, uint64)
}
