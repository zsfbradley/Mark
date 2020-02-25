package mFace

type MError interface {
	error
	NotNil() bool
	TCError() []byte
}
