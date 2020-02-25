package mNet

import (
	"log"
	"markV5/mConst"
	"markV5/mError"
	"markV5/mFace"
	"net"
)

func DefaultServer() (mFace.MServer, mFace.MError) {
	config := defaultConfig()
	return ServerWithConfig(config)
}

func ServerWithConfig(config mFace.MConfig) (mFace.MServer, mFace.MError) {
	s := &server{
		config:   config,
		status:   mConst.MServeStatus_UnStart,
		cm:       newConnManager(),
		mm:       newMsgManager(),
		rm:       newRouteManager(),
		listener: nil,
		eFuncs:   make([]mFace.MEntranceFunc, 0),
	}

	s.cm.SuperiorServer(s)
	s.mm.SuperiorServer(s)
	s.rm.SuperiorServer(s)

	if err := s.Load(); err.NotNil() {
		return nil, err
	}

	return s, mError.Nil
}

// server 的主要任务是负责连接connManager、msgManager、routeManager，承担中心轴的角色
type server struct {
	config mFace.MConfig
	status mConst.MServeStatus

	cm mFace.MConnManager
	mm mFace.MMsgManager
	rm mFace.MRouteManager

	listener *net.Listener
	eFuncs   []mFace.MEntranceFunc
}

func (s *server) RegisterEntranceFunc(eFunc mFace.MEntranceFunc) mFace.MError {
	if eFunc == nil {
		return mError.NilParam
	}

	s.eFuncs = append(s.eFuncs, eFunc)

	return mError.Nil
}

func (s *server) Load() mFace.MError {
	if err := s.config.Load(); err.NotNil() {
		return err
	}

	log.Printf("[%s] Server Load", s.config.Name())

	if err := s.cm.Load(); err.NotNil() {
		return err
	}

	if err := s.mm.Load(); err.NotNil() {
		return err
	}

	if err := s.rm.Load(); err.NotNil() {
		return err
	}

	s.status = mConst.MServeStatus_Load

	address := net.JoinHostPort(s.config.Host(), s.config.Port())
	l, err := net.Listen(s.config.NetWork(), address)
	if err != nil {
		return mError.SystemError(err)
	}

	s.listener = &l

	log.Printf("[%s] Server Load Success", s.config.Name())
	return mError.Nil
}

func (s *server) Start() mFace.MError {
	log.Printf("[%s] Server Start", s.config.Name())

	if err := s.cm.Start(); err.NotNil() {
		return err
	}

	if err := s.mm.Start(); err.NotNil() {
		return err
	}

	if err := s.rm.Start(); err.NotNil() {
		return err
	}

	s.status = mConst.MServeStatus_Start

	go s.acceptConn()

	for _, eFunc := range s.eFuncs {
		if err := eFunc(); err.NotNil() {
			return err
		}
	}

	log.Printf("[%s] Server Start Success", s.config.Name())
	return mError.Nil
}

func (s *server) Reload() mFace.MError {
	log.Printf("[%s] Server Reload", s.config.Name())

	if err := s.config.Reload(); err.NotNil() {
		return err
	}

	if err := s.cm.Reload(); err.NotNil() {
		return err
	}

	if err := s.mm.Reload(); err.NotNil() {
		return err
	}

	if err := s.rm.Reload(); err.NotNil() {
		return err
	}

	log.Printf("[%s] Server Reload Success", s.config.Name())
	s.status = mConst.MServeStatus_Reload
	return mError.Nil
}

func (s *server) Stop() mFace.MError {
	log.Printf("[%s] Server Stop", s.config.Name())

	s.status = mConst.MServeStatus_StartEnding

	if err := s.cm.StartEnding(); err.NotNil() {
		return err
	}

	if err := s.mm.StartEnding(); err.NotNil() {
		return err
	}

	if err := s.rm.StartEnding(); err.NotNil() {
		return err
	}

	if err := s.rm.OfficialEnding(); err.NotNil() {
		return err
	}

	if err := s.mm.OfficialEnding(); err.NotNil() {
		return err
	}

	if err := s.cm.OfficialEnding(); err.NotNil() {
		return err
	}

	s.status = mConst.MServeStatus_OfficialEnding
	listener := *(s.listener)
	if err := listener.Close(); err != nil {
		return mError.SystemError(err)
	}

	log.Printf("[%s] Server Stop Success", s.config.Name())
	s.status = mConst.MServeStatus_Stopped
	return mError.Nil
}

func (s *server) Status() mConst.MServeStatus {
	return s.status
}

func (s *server) Config() mFace.MConfig {
	return s.config
}

func (s *server) ConnManager() mFace.MConnManager {
	return s.cm
}

func (s *server) MsgManager() mFace.MMsgManager {
	return s.mm
}

func (s *server) RouteManager() mFace.MRouteManager {
	return s.rm
}

func (s *server) RegisterRoute(routeID string, handleFunc mFace.RouteHandleFunc) mFace.MError {
	if handleFunc == nil {
		return mError.NilParam
	}
	newRouteHandler := newRouteHandler(routeID, handleFunc)
	return s.rm.RegisterNewRoute(newRouteHandler)
}

func (s *server) RegisterRoutes(routes map[string]mFace.RouteHandleFunc) mFace.MError {
	if routes == nil || len(routes) == 0 {
		return mError.NilParam
	}
	for routeID, handleFunc := range routes {
		newRouteHandler := newRouteHandler(routeID, handleFunc)
		err := s.rm.RegisterNewRoute(newRouteHandler)
		if err.NotNil() {
			return err
		}
	}
	return mError.Nil
}

func (s *server) RegisterFilter(msgFilterType mFace.MsgFilterType, msgFilterFunc mFace.MsgFilterFunc) mFace.MError {
	if msgFilterFunc == nil {
		return mError.NilParam
	}
	newFilter := newMsgFilter(msgFilterType, msgFilterFunc)
	s.mm.RegisterNewFilter(newFilter)
	return mError.Nil
}

// - private methods

func (s *server) acceptConn() {
	log.Printf("[%s] Server Start Accept Conn", s.config.Name())
	listener := *(s.listener)
	for {
		conn, err := listener.Accept()
		if err != nil {
			if s.status >= mConst.MServeStatus_StartEnding { // 正式关停服务
				break
			}
			continue
		}

		if err := s.cm.AcceptNewConn(&conn); err.NotNil() {
			log.Println(err.Error())
		}
	}
	log.Printf("[%s] Server Stop Accept Conn", s.config.Name())
}
