package mNet

import (
	"log"
	"markV5/mConst"
	"markV5/mError"
	"markV5/mFace"
	"sync"
)

func newMsgManager() mFace.MMsgManager {
	mm := &msgManager{
		status:      mConst.MServeStatus_UnStart,
		ss:          nil,
		filters:     make([]mFace.MMsgFilter, 0),
		filtersLock: sync.RWMutex{},
		msgInChan:   newChannel(),
		msgOutChan:  newChannel(),
	}
	return mm
}

// msgManager 的责任是负责有限制且有序的传递消息，并通过外部注册的过滤器过滤消息
type msgManager struct {
	status mConst.MServeStatus
	ss     mFace.MServer

	filters     []mFace.MMsgFilter
	filtersLock sync.RWMutex

	msgInChan  mFace.MChannel
	msgOutChan mFace.MChannel
}

func (mm *msgManager) SuperiorServer(server mFace.MServer) {
	mm.ss = server
}

func (mm *msgManager) Load() mFace.MError {
	log.Printf("MsgManager Load")
	mm.status = mConst.MServeStatus_Load
	mm.msgInChan.SetSize(mm.ss.Config().MsgManagerInChanSize())
	if err := mm.msgInChan.Load(); err != nil {
		return err
	}
	mm.msgOutChan.SetSize(mm.ss.Config().MsgManagerOutChanSize())
	if err := mm.msgOutChan.Load(); err != nil {
		return err
	}
	return mError.Nil
}

func (mm *msgManager) Start() mFace.MError {
	//log.Println(mm.filters)
	//for _ , filter := range mm.filters {
	//	log.Println(filter)
	//	log.Println(filter.Type())
	//	log.Println(filter.FilterFunc())
	//}

	log.Printf("MsgManager Start")
	mm.status = mConst.MServeStatus_Start
	if err := mm.msgInChan.Start(); err != nil {
		return err
	}
	if err := mm.msgOutChan.Start(); err != nil {
		return err
	}

	go mm.acceptResponse()
	go mm.acceptRequest()

	return mError.Nil
}

func (mm *msgManager) Reload() mFace.MError {
	log.Printf("MsgManager Reload")
	mm.status = mConst.MServeStatus_Reload
	mm.msgInChan.SetSize(mm.ss.Config().MsgManagerInChanSize())
	if err := mm.msgInChan.Reload(); err != nil {
		return err
	}
	mm.msgOutChan.SetSize(mm.ss.Config().MsgManagerOutChanSize())
	if err := mm.msgOutChan.Reload(); err != nil {
		return err
	}
	return mError.Nil
}

func (mm *msgManager) StartEnding() mFace.MError {
	log.Printf("MsgManager Start Ending")
	mm.status = mConst.MServeStatus_StartEnding
	if err := mm.msgInChan.StartEnding(); err != nil {
		return err
	}
	if err := mm.msgOutChan.StartEnding(); err != nil {
		return err
	}
	return mError.Nil
}

func (mm *msgManager) OfficialEnding() mFace.MError {
	log.Printf("MsgManager Official Ending")
	mm.status = mConst.MServeStatus_OfficialEnding
	if err := mm.msgInChan.OfficialEnding(); err != nil {
		return err
	}
	if err := mm.msgOutChan.OfficialEnding(); err != nil {
		return err
	}
	return mError.Nil
}

func (mm *msgManager) Status() mConst.MServeStatus {
	return mm.status
}

func (mm *msgManager) DataInChannel() mFace.MChannel {
	return mm.msgInChan
}

func (mm *msgManager) DataOutChannel() mFace.MChannel {
	return mm.msgOutChan
}

// 注册消息过滤器
func (mm *msgManager) RegisterNewFilter(filter mFace.MMsgFilter) {
	mm.filtersLock.Lock()
	defer mm.filtersLock.Unlock()

	mm.filters = append(mm.filters, filter)
}

// - private methods

func (mm *msgManager) acceptResponse() {
	for {
		newData, ok := mm.msgOutChan.Out() // 取出
		if !ok {
			if mm.msgOutChan.Status() >= mConst.MServeStatus_OfficialEnding {
				break
			} else {
				log.Println("Msg Manager MsgOutChan !ok")
				continue
			}
		}
		newMsg := newData.(mFace.MMessage) // 断言成 mFace.MMessage类型
		go mm.filterMsg(mFace.MsgFilterType_Response, mm.ss.ConnManager().DataOutChannel(), newMsg)
	}
	log.Println("Msg Manager MsgOutChan End Work")
}

func (mm *msgManager) acceptRequest() {
	for {
		newData, ok := mm.msgInChan.Out() // 取出
		if !ok {
			if mm.msgInChan.Status() >= mConst.MServeStatus_OfficialEnding {
				break
			} else {
				log.Println("Msg Manager MsgInChan !ok")
				continue
			}
		}
		newMsg := newData.(mFace.MMessage) // 断言成 mFace.MMessage类型
		go mm.filterMsg(mFace.MsgFilterType_Request, mm.ss.RouteManager().DataInChannel(), newMsg)
	}
	log.Println("Msg Manager MsgInChan End Work")
}

func (mm *msgManager) filterMsg(msgType mFace.MsgFilterType, finalChannel mFace.MChannel, msg mFace.MMessage) {
	mm.filtersLock.RLock()
	defer mm.filtersLock.RUnlock()

	filterResult := true
	for _, filter := range mm.filters {
		if filter.Type() != msgType {
			continue
		}
		filterFunc := filter.FilterFunc()
		if pass := filterFunc(msg); !pass {
			filterResult = false
			break
		}
	}

	if filterResult {
		// 塞进 finalChannel 中
		err := finalChannel.In(msg)
		if err.NotNil() {
			log.Println(err)
		}
		return
	}

	log.Printf("Msg %v not pass filter", msg)
}
