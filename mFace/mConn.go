package mFace

type MConn interface {
	MServeLoad
	MServeStart
	MServeStatus
	MServeEnding

	SuperiorManager(MConnManager)
	ConnID() string
	Response(MMessage)
}
