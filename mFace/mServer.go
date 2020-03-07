package mFace

type EntranceFunc func() error

type MServer interface {
	MServeBase
	MServeStop

	RegisterEntranceFunc(EntranceFunc) error
}
