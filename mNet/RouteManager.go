package mNet

import (
	"../mConst"
	"../mFace"
	"log"
)

func newRouteManager() mFace.MRouteManager {
	rm := &routeManager{
		server: nil,
		status: mConst.MServe_Status_UnStart,
	}
	return rm
}

type routeManager struct {
	server mFace.MServer
	status mConst.MServe_Status
}

func (rm *routeManager) BindServer(s mFace.MServer) {
	rm.server = s
}

func (rm *routeManager) Status() mConst.MServe_Status {
	return rm.status
}

func (rm *routeManager) Load() error {
	rm.status = mConst.MServe_Status_Load
	log.Printf("[%s] routeManager start load", rm.server.Config().Name)

	log.Printf("[%s] routeManager finish load", rm.server.Config().Name)

	return nil
}

func (rm *routeManager) Start() error {
	rm.status = mConst.MServe_Status_Start
	log.Printf("[%s] routeManager starting", rm.server.Config().Name)

	log.Printf("[%s] routeManager started", rm.server.Config().Name)
	return nil
}

func (rm *routeManager) Reload() error {
	rm.status = mConst.MServe_Status_Reload
	log.Printf("[%s] routeManager start reload", rm.server.Config().Name)

	log.Printf("[%s] routeManager finish reload", rm.server.Config().Name)
	rm.status = mConst.MServe_Status_Start
	return nil
}

func (rm *routeManager) StartEnding() error {
	rm.status = mConst.MServe_Status_StartEnding
	log.Printf("[%s] routeManager start ending", rm.server.Config().Name)
	return nil
}

func (rm *routeManager) OfficialEnding() error {
	rm.status = mConst.MServe_Status_OfficialEnding
	log.Printf("[%s] routeManager official ending", rm.server.Config().Name)

	rm.status = mConst.MServe_Status_Stoped
	log.Printf("[%s] connManager stoped", rm.server.Config().Name)
	return nil
}

// ----------------------------------------------------------- private methods
