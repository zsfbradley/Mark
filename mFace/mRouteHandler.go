package mFace

type RouteHandleFunc func(MMessage, MMessage) MError

type MRouteHandler interface {
	RouteID() string
	RouteHandler() RouteHandleFunc
}
