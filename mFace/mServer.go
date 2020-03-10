package mFace

type MServer interface {
	MServe
	MServeLoad
	MServeStop

	Config() ServerConfig
}
