package mFace

type MConfig interface {
	MServeLoad
	Name() string
	NetWork() string
	Host() string
	Port() string
	RouteManagerInChanSize() uint64
	RouteManagerOutChanSize() uint64
	MsgManagerInChanSize() uint64
	MsgManagerOutChanSize() uint64
	ConnManagerInChanSize() uint64
	ConnManagerOutChanSize() uint64
}
