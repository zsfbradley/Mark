package mNet

import (
	"../mConst"
	"../mFace"
	"log"
)

func newConnManager() mFace.MConnManager {
	cm := &connManager{
		status: mConst.MServe_Status_UnStart,
	}
	return cm
}

type connManager struct {
	server mFace.MServer
	status mConst.MServe_Status
}

func (cm *connManager) BindServer(s mFace.MServer) {
	cm.server = s
}

func (cm *connManager) Status() mConst.MServe_Status {
	return cm.status
}

func (cm *connManager) Load() error {
	cm.status = mConst.MServe_Status_Load
	log.Printf("[%s] connManager start load" , cm.server.Config().Name)

	log.Printf("[%s] connManager finish load" , cm.server.Config().Name)

	return nil
}

func (cm *connManager) Start() error {
	cm.status = mConst.MServe_Status_Start
	log.Printf("[%s] connManager starting" , cm.server.Config().Name)

	log.Printf("[%s] connManager started" , cm.server.Config().Name)

	return nil
}

func (cm *connManager) Reload() error {
	cm.status = mConst.MServe_Status_Reload
	log.Printf("[%s] connManager start reload" , cm.server.Config().Name)

	log.Printf("[%s] connManager finish reload" , cm.server.Config().Name)
	cm.status = mConst.MServe_Status_Start

	return nil
}

func (cm *connManager) StartEnding() error {
	cm.status = mConst.MServe_Status_StartEnding
	log.Printf("[%s] connManager start ending" , cm.server.Config().Name)
	return nil
}

func (cm *connManager) OfficialEnding() error {
	cm.status = mConst.MServe_Status_OfficialEnding
	log.Printf("[%s] connManager official ending" , cm.server.Config().Name)

	cm.status = mConst.MServe_Status_Stoped

	return nil
}

// ----------------------------------------------------------- private methods