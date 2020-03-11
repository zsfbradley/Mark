package mFace

type MConn interface {
	MServe
	MServeLoad
	MServeEnding

	BindManager(MConnManager)
	ID() string
	ReplyResponse([]byte) error
}
