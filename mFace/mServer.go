package mFace

type EntranceFunc func() error

type MServer interface {
	MServe
	MServeLoad
	MServeStop

	Config() ServerConfig
	RegisterEncFunc(EntranceFunc) error
	RegisterCodecCreator(CodecCreatorFunc) error
}
