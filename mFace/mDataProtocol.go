package mFace

type MDataProtocol interface {
	Unmarshal([]byte)
	Marshal(...[]byte) []byte
	CompletedDataChannel() chan [][]byte
}
