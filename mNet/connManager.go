package mNet

import (
	"../mConst"
	"../mFace"
	"fmt"
	"log"
	"sync"
)

func newConnManager() mFace.MConnManager {
	cm := &connManager{
		server:      nil,
		status:      mConst.MServe_Status_UnStart,
		maxConn:     0,
		connMap:     make(map[string]mFace.MConn),
		connMapLock: sync.RWMutex{},
	}
	return cm
}

type connManager struct {
	server mFace.MServer
	status mConst.MServe_Status

	maxConn     int64
	connMap     map[string]mFace.MConn
	connMapLock sync.RWMutex
}

func (cm *connManager) BindServer(s mFace.MServer) {
	cm.server = s
}

func (cm *connManager) Status() mConst.MServe_Status {
	return cm.status
}

func (cm *connManager) Load() error {
	cm.status = mConst.MServe_Status_Load
	log.Printf("[%s] connManager start load", cm.server.Config().Name)

	cm.maxConn = cm.server.Config().MaxConn

	log.Printf("[%s] connManager finish load", cm.server.Config().Name)

	return nil
}

func (cm *connManager) Start() error {
	cm.status = mConst.MServe_Status_Start
	log.Printf("[%s] connManager starting", cm.server.Config().Name)

	log.Printf("[%s] connManager started", cm.server.Config().Name)

	return nil
}

func (cm *connManager) Reload() error {
	cm.status = mConst.MServe_Status_Reload
	log.Printf("[%s] connManager start reload", cm.server.Config().Name)

	log.Printf("[%s] connManager finish reload", cm.server.Config().Name)
	cm.status = mConst.MServe_Status_Start

	return nil
}

func (cm *connManager) StartEnding() error {
	cm.status = mConst.MServe_Status_StartEnding
	log.Printf("[%s] connManager start ending", cm.server.Config().Name)
	return nil
}

func (cm *connManager) OfficialEnding() error {
	cm.status = mConst.MServe_Status_OfficialEnding
	log.Printf("[%s] connManager official ending", cm.server.Config().Name)

	cm.status = mConst.MServe_Status_Stoped
	log.Printf("[%s] connManager stoped", cm.server.Config().Name)

	return nil
}

func (cm *connManager) AcceptNewConn(codec mFace.MCodec) {
	cm.connMapLock.Lock()
	defer cm.connMapLock.Unlock()

	if !cm.checkMaxConnLimit() {
		cm.refuseConn(codec)
		return
	}

	conn := newConn(DefaultIDG().EncNewID(), codec)
	cm.connMap[conn.ID()] = conn
	conn.BindManager(cm)
	log.Printf("[%s] conn load : %v", conn.ID(), conn.Load())
	log.Printf("[%s] conn start : %v", conn.ID(), conn.Start())
}

func (cm *connManager) DeleteConn(id string) {
	cm.connMapLock.Lock()
	defer cm.connMapLock.Unlock()

	delete(cm.connMap, id)
}

// ----------------------------------------------------------- private methods

func (cm *connManager) checkMaxConnLimit() bool {
	if cm.maxConn == 0 {
		return true
	}
	currentLen := int64(len(cm.connMap))
	return currentLen < cm.maxConn
}

func (cm *connManager) refuseConn(codec mFace.MCodec) {
	_ = codec.WriteResponse([]byte(fmt.Sprintf(mConst.Refuse_Connect_Of_Limit, cm.server.Config().Name)))
	_ = codec.Close()
}
