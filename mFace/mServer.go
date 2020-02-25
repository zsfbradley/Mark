package mFace

type MEntranceFunc func() MError

type MServer interface {
	MServeLoad
	MServeStart
	MServeStop
	MServeStatus

	Config() MConfig

	ConnManager() MConnManager
	MsgManager() MMsgManager
	RouteManager() MRouteManager

	RegisterRoute(string, RouteHandleFunc) MError
	RegisterRoutes(map[string]RouteHandleFunc) MError
	RegisterFilter(MsgFilterType, MsgFilterFunc) MError
	RegisterEntranceFunc(MEntranceFunc) MError
}
