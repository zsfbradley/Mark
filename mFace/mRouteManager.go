package mFace

type MRouteManager interface {
	MServe
	MServeLoad
	MServeEnding

	BindServer(MServer)
}
