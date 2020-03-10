package mFace

type MConnManager interface {
	MServe
	MServeLoad
	MServeEnding

	BindServer(MServer)
}
