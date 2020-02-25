package mNet

import (
	"log"
	"markV5/mConst"
	"markV5/mError"
	"markV5/mFace"
	"markV5/mTool"
	"net"
	"sync"
)

func newConnManager() mFace.MConnManager {
	cm := &connManager{
		status:     mConst.MServeStatus_UnStart,
		ss:         nil,
		conns:      make(map[string]mFace.MConn),
		connsLock:  sync.RWMutex{},
		msgInChan:  newChannel(),
		msgOutChan: newChannel(),
	}
	return cm
}

// connManager 的主要责任是负责有限制的接收新的外部链接，并负责其完整的生命周期
type connManager struct {
	status mConst.MServeStatus
	ss     mFace.MServer

	conns     map[string]mFace.MConn
	connsLock sync.RWMutex

	msgInChan  mFace.MChannel
	msgOutChan mFace.MChannel
}

func (cm *connManager) SuperiorServer(server mFace.MServer) {
	cm.ss = server
}

func (cm *connManager) Load() mFace.MError {
	log.Printf("ConnManager Load")
	cm.status = mConst.MServeStatus_Load
	cm.msgInChan.SetSize(cm.ss.Config().ConnManagerInChanSize())
	if err := cm.msgInChan.Load(); err != nil {
		return err
	}
	cm.msgOutChan.SetSize(cm.ss.Config().ConnManagerOutChanSize())
	if err := cm.msgOutChan.Load(); err != nil {
		return err
	}
	return mError.Nil
}

func (cm *connManager) Start() mFace.MError {
	log.Printf("ConnManager Start")
	cm.status = mConst.MServeStatus_Start
	if err := cm.msgInChan.Start(); err != nil {
		return err
	}
	if err := cm.msgOutChan.Start(); err != nil {
		return err
	}

	go cm.acceptResponse()
	go cm.acceptRequest()

	return mError.Nil
}

func (cm *connManager) Reload() mFace.MError {
	log.Printf("ConnManager Reload")
	cm.status = mConst.MServeStatus_Reload
	cm.msgInChan.SetSize(cm.ss.Config().ConnManagerInChanSize())
	if err := cm.msgInChan.Reload(); err != nil {
		return err
	}
	cm.msgOutChan.SetSize(cm.ss.Config().ConnManagerOutChanSize())
	if err := cm.msgOutChan.Reload(); err != nil {
		return err
	}

	cm.connsLock.RLock()
	defer cm.connsLock.RUnlock()

	for _, conn := range cm.conns {
		if err := conn.Reload(); err.NotNil() {
			return err
		}
	}

	return mError.Nil
}

func (cm *connManager) StartEnding() mFace.MError {
	log.Printf("ConnManager Start Ending")
	cm.status = mConst.MServeStatus_StartEnding
	if err := cm.msgInChan.StartEnding(); err != nil {
		return err
	}
	if err := cm.msgOutChan.StartEnding(); err != nil {
		return err
	}

	cm.connsLock.RLock()
	defer cm.connsLock.RUnlock()

	for _, conn := range cm.conns {
		if err := conn.StartEnding(); err.NotNil() {
			return err
		}
	}

	return mError.Nil
}

func (cm *connManager) OfficialEnding() mFace.MError {
	log.Printf("ConnManager Official Ending")
	cm.status = mConst.MServeStatus_OfficialEnding
	if err := cm.msgInChan.OfficialEnding(); err != nil {
		return err
	}
	if err := cm.msgOutChan.OfficialEnding(); err != nil {
		return err
	}

	cm.connsLock.RLock()
	defer cm.connsLock.RUnlock()

	for id, conn := range cm.conns {
		if err := conn.OfficialEnding(); err.NotNil() {
			return err
		}
		delete(cm.conns, id)
	}

	return mError.Nil
}

func (cm *connManager) Status() mConst.MServeStatus {
	return cm.status
}

func (cm *connManager) DataInChannel() mFace.MChannel {
	return cm.msgInChan
}

func (cm *connManager) DataOutChannel() mFace.MChannel {
	return cm.msgOutChan
}

func (cm *connManager) AcceptNewConn(conn *net.Conn) mFace.MError {
	if conn == nil {
		return mError.NilParam
	}

	cm.connsLock.Lock()
	defer cm.connsLock.Unlock()

	randomID := mTool.GetRandString(32)
	newConn := newConn(randomID, conn)
	newConn.SuperiorManager(cm) // 子 Conn 持有父级

	cm.conns[randomID] = newConn

	if err := newConn.Load(); err.NotNil() {
		return err
	}

	if err := newConn.Start(); err.NotNil() {
		return err
	}

	return mError.Nil
}

func (cm *connManager) ConnOut(connID string) mFace.MError {
	if connID == "" {
		return mError.NilParam
	}

	cm.connsLock.Lock()
	defer cm.connsLock.Unlock()

	delete(cm.conns, connID)

	return mError.Nil
}

// - private methods

func (cm *connManager) acceptRequest() {
	for {
		newData, ok := cm.msgInChan.Out()
		if !ok {
			if cm.msgInChan.Status() >= mConst.MServeStatus_OfficialEnding {
				break
			} else {
				log.Println("Conn Manager MsgInChan !ok")
				continue
			}
		}
		newMsg := newData.(mFace.MMessage) // 断言成 mFace.MMessage类型
		if err := cm.ss.MsgManager().DataInChannel().In(newMsg); err.NotNil() {
			log.Println(err.Error())
		}
	}
	log.Println("Conn Manager MsgInChan End Work")
}

func (cm *connManager) acceptResponse() {
	for {
		newData, ok := cm.msgOutChan.Out() // 取出
		if !ok {
			if cm.msgOutChan.Status() >= mConst.MServeStatus_OfficialEnding {
				break
			} else {
				log.Println("Conn Manager MsgOutChan !ok")
				continue
			}
		}
		newMsg := newData.(mFace.MMessage) // 断言成 mFace.MMessage类型
		go cm.notifyResponse(newMsg)

	}
	log.Println("Conn Manager MsgOutChan End Work")
}

func (cm *connManager) notifyResponse(rsp mFace.MMessage) {
	conn, err := cm.getConn(rsp.ConnID())
	if err.NotNil() {
		log.Println(err.Error())
		return
	}

	conn.Response(rsp)
}

func (cm *connManager) getConn(id string) (mFace.MConn, mFace.MError) {
	if id == "" {
		return nil, mError.ErrorParam
	}
	cm.connsLock.RLock()
	defer cm.connsLock.RUnlock()

	conn, exist := cm.conns[id]
	if !exist {
		return nil, mError.ConnIDUnExist
	}

	return conn, mError.Nil
}
