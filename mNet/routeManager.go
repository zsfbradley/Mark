package mNet

import (
	"log"
	"markV5/mConst"
	"markV5/mError"
	"markV5/mFace"
	"sync"
)

func newRouteManager() mFace.MRouteManager {
	rm := &routeManager{
		status:     mConst.MServeStatus_UnStart,
		ss:         nil,
		routes:     make(map[string]mFace.MRouteHandler),
		routesLock: sync.RWMutex{},
		hooks:      make([]mFace.MRouterHookFunc, 0),
		hooksLock:  sync.RWMutex{},
		msgInChan:  newChannel(),
		msgOutChan: newChannel(),
	}
	return rm
}

// routeManager 的责任是负责有限制的注册路由，并将响应消息回传
type routeManager struct {
	status mConst.MServeStatus
	ss     mFace.MServer

	routes     map[string]mFace.MRouteHandler
	routesLock sync.RWMutex

	hooks     []mFace.MRouterHookFunc
	hooksLock sync.RWMutex

	msgInChan  mFace.MChannel
	msgOutChan mFace.MChannel
}

func (rm *routeManager) SuperiorServer(server mFace.MServer) {
	rm.ss = server
}

func (rm *routeManager) Load() mFace.MError {
	log.Printf("RouteManager Load")
	rm.status = mConst.MServeStatus_Load
	rm.msgInChan.SetSize(rm.ss.Config().RouteManagerInChanSize())
	if err := rm.msgInChan.Load(); err != nil {
		return err
	}
	rm.msgOutChan.SetSize(rm.ss.Config().RouteManagerOutChanSize())
	if err := rm.msgOutChan.Load(); err != nil {
		return err
	}
	return mError.Nil
}

func (rm *routeManager) Start() mFace.MError {
	log.Printf("RouteManager Start")
	rm.status = mConst.MServeStatus_Start
	if err := rm.msgInChan.Start(); err != nil {
		return err
	}
	if err := rm.msgOutChan.Start(); err != nil {
		return err
	}

	go rm.acceptResponse()
	go rm.acceptRequest()

	return mError.Nil
}

func (rm *routeManager) Reload() mFace.MError {
	log.Printf("RouteManager Reload")
	rm.status = mConst.MServeStatus_Reload
	rm.msgInChan.SetSize(rm.ss.Config().RouteManagerInChanSize())
	if err := rm.msgInChan.Reload(); err != nil {
		return err
	}
	rm.msgOutChan.SetSize(rm.ss.Config().RouteManagerOutChanSize())
	if err := rm.msgOutChan.Reload(); err != nil {
		return err
	}
	return mError.Nil
}

func (rm *routeManager) StartEnding() mFace.MError {
	log.Printf("RouteManager Start Ending")
	rm.status = mConst.MServeStatus_StartEnding
	if err := rm.msgInChan.StartEnding(); err != nil {
		return err
	}
	if err := rm.msgOutChan.StartEnding(); err != nil {
		return err
	}
	return mError.Nil
}

func (rm *routeManager) OfficialEnding() mFace.MError {
	log.Printf("RouteManager Official Ending")
	rm.status = mConst.MServeStatus_OfficialEnding
	if err := rm.msgInChan.OfficialEnding(); err != nil {
		return err
	}
	if err := rm.msgOutChan.OfficialEnding(); err != nil {
		return err
	}
	return mError.Nil
}

func (rm *routeManager) Status() mConst.MServeStatus {
	return rm.status
}

func (rm *routeManager) DataInChannel() mFace.MChannel {
	return rm.msgInChan
}

func (rm *routeManager) DataOutChannel() mFace.MChannel {
	return rm.msgOutChan
}

// 注册路由
func (rm *routeManager) RegisterNewRoute(newRoute mFace.MRouteHandler) mFace.MError {
	rm.routesLock.Lock()
	defer rm.routesLock.Unlock()

	if _, exist := rm.routes[newRoute.RouteID()]; exist {
		return mError.RouteExist
	}

	rm.routes[newRoute.RouteID()] = newRoute

	return mError.Nil
}

// 取消注册路由
func (rm *routeManager) UnRegisterRoute(routeID string) mFace.MError {
	rm.routesLock.Lock()
	defer rm.routesLock.Unlock()

	if _, exist := rm.routes[routeID]; !exist {
		return mError.RouteUnExist
	}

	delete(rm.routes, routeID)

	return mError.Nil
}

func (rm *routeManager) AddHook(hookFunc mFace.MRouterHookFunc) {
	rm.hooksLock.Lock()
	defer rm.hooksLock.Unlock()

	rm.hooks = append(rm.hooks, hookFunc)
}

// - private methods

func (rm *routeManager) handler(routeID string) (mFace.MRouteHandler, mFace.MError) {
	rm.routesLock.RLock()
	defer rm.routesLock.RUnlock()

	handler, exist := rm.routes[routeID]
	if !exist {
		return nil, mError.RouteUnExist
	}

	newHandler := handler

	return newHandler, mError.Nil
}

func (rm *routeManager) hook(handler mFace.RouteHandleFunc) mFace.RouteHandleFunc {
	rm.hooksLock.RLock()
	defer rm.hooksLock.RUnlock()

	for _, hookFunc := range rm.hooks {
		handler = hookFunc(handler)
	}
	return handler
}

// 接纳新响应消息
func (rm *routeManager) acceptResponse() {
	for {
		newData, ok := rm.msgOutChan.Out()
		if !ok {
			if rm.msgOutChan.Status() >= mConst.MServeStatus_OfficialEnding {
				break
			} else {
				log.Println("Route Manager MsgOutChan !ok")
				continue
			}
		}
		newMsg := newData.(mFace.MMessage) // 断言成 mFace.MMessage类型
		if err := rm.ss.MsgManager().DataOutChannel().In(newMsg); err.NotNil() {
			log.Println(err)
		}
	}
	log.Println("Route Manager MsgOutChan End Work")
}

// 接纳新消息
func (rm *routeManager) acceptRequest() {
	for {
		newData, ok := rm.msgInChan.Out()
		if !ok {
			if rm.msgInChan.Status() >= mConst.MServeStatus_OfficialEnding {
				break
			} else {
				log.Println("Route Manager MsgInChan !ok")
				continue
			}
		}
		newMsg := newData.(mFace.MMessage) // 断言成 mFace.MMessage类型
		go rm.handleNewMsg(newMsg)
	}
	log.Println("Route Manager MsgInChan End Work")
}

// 处理新消息
func (rm *routeManager) handleNewMsg(newMsg mFace.MMessage) {
	handler, err := rm.handler(newMsg.RouteID())
	if err.NotNil() {
		log.Println(err)
		return
	}
	handleFunc := handler.RouteHandler()
	rsp := newMessage(newMsg.RouteID(), newMsg.ConnID(), []byte{})
	handleFunc = rm.hook(handleFunc)
	err = handleFunc(newMsg, rsp)
	if err.NotNil() {
		rsp.Unmarshal(err.TCError())
		log.Println(err.Error())
	}
	if err := rm.msgOutChan.In(rsp); err.NotNil() {
		log.Println(err.Error())
	}
}
