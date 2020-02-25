package mFace

type MRouteManager interface {
	MServeLoad
	MServeStart
	MServeEnding
	MServeStatus
	MServeDataChannel

	SuperiorServer(MServer)
	RegisterNewRoute(MRouteHandler) MError
	UnRegisterRoute(string) MError
	AddHook(MRouterHookFunc)
}
