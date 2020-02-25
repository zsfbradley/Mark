package mFace

type MChannel interface {
	MServeLoad
	MServeStart
	MServeStatus
	MServeEnding

	SetSize(uint64)
	In(interface{}) MError
	Out() (interface{}, bool)
	Close() MError
}
