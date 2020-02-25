package mFace

type MMsgManager interface {
	MServeLoad
	MServeStart
	MServeEnding
	MServeStatus
	MServeDataChannel

	SuperiorServer(MServer)
	RegisterNewFilter(MMsgFilter)
}
