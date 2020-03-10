package mFace

type MMessageManager interface {
	MServe
	MServeLoad
	MServeEnding

	BindServer(MServer)
}
