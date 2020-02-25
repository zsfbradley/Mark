package mNet

import "markV5/mFace"

func newRouteHandler(routeID string, handleFunc mFace.RouteHandleFunc) mFace.MRouteHandler {
	rh := &routeHandler{
		routeID:         routeID,
		routeHandleFunc: handleFunc,
	}
	return rh
}

type routeHandler struct {
	routeID         string
	routeHandleFunc mFace.RouteHandleFunc
}

func (rh *routeHandler) RouteID() string {
	return rh.routeID
}

func (rh *routeHandler) RouteHandler() mFace.RouteHandleFunc {
	return rh.routeHandleFunc
}
