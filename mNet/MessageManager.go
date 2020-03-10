package mNet

import (
	"../mConst"
	"../mFace"
	"log"
)

func newMessageManager() mFace.MMessageManager {
	mm := &messageManager{
		status: mConst.MServe_Status_UnStart,
	}
	return mm
}

type messageManager struct {
	server mFace.MServer
	status mConst.MServe_Status
}

func (mm *messageManager) BindServer(s mFace.MServer) {
	mm.server = s
}

func (mm *messageManager) Status() mConst.MServe_Status {
	return mm.status
}

func (mm *messageManager) Load() error {
	mm.status = mConst.MServe_Status_Load
	log.Printf("[%s] messageManager start load" , mm.server.Config().Name)

	log.Printf("[%s] messageManager finish load" , mm.server.Config().Name)

	return nil
}

func (mm *messageManager) Start() error {
	mm.status = mConst.MServe_Status_Start
	log.Printf("[%s] messageManager starting" , mm.server.Config().Name)

	log.Printf("[%s] messageManager started" , mm.server.Config().Name)
	return nil
}

func (mm *messageManager) Reload() error {
	mm.status = mConst.MServe_Status_Reload
	log.Printf("[%s] messageManager start reload" , mm.server.Config().Name)

	log.Printf("[%s] messageManager finish reload" , mm.server.Config().Name)
	mm.status = mConst.MServe_Status_Start
	return nil
}

func (mm *messageManager) StartEnding() error {
	mm.status = mConst.MServe_Status_StartEnding
	log.Printf("[%s] messageManager start ending" , mm.server.Config().Name)
	return nil
}

func (mm *messageManager) OfficialEnding() error {
	mm.status = mConst.MServe_Status_OfficialEnding
	log.Printf("[%s] messageManager official ending" , mm.server.Config().Name)

	mm.status = mConst.MServe_Status_Stoped
	return nil
}

// ----------------------------------------------------------- private methods