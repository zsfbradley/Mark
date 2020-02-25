package mFace

import "net"

type MConnManager interface {
	MServeLoad
	MServeStart
	MServeEnding
	MServeStatus
	MServeDataChannel

	SuperiorServer(MServer)
	AcceptNewConn(*net.Conn) MError
	ConnOut(string) MError
}
