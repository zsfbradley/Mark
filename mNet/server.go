package mNet

import (
	"../mConst"
	"../mFace"
	"log"
	"net"
)

func NewServer(sConfig ServerConfig) (mFace.MServer, error) {
	sConfig.parse()

	s := &server{
		config:   sConfig,
		status:   mConst.MServe_Status_UnStart,
		listener: nil,
	}

	err := s.Load()

	return s, err
}

type ServerConfig struct {
	Name    string // option , default is "default"
	Host    string // option , default is 0.0.0.0
	Port    string // require	, server not gonna listen without port
	Network string // option , default is tcp
}

func (sc *ServerConfig) parse() {
	if sc.Name == "" {
		sc.Name = mConst.Framework_Name + mConst.Framework_Version
	}
	if sc.Host == "" {
		sc.Host = mConst.Default_Host
	}
	if sc.Network == "" {
		sc.Network = mConst.Network_TCP
	}
}

type server struct {
	config   ServerConfig
	status   mConst.MServe_Status
	listener net.Listener
}

func (s *server) Status() mConst.MServe_Status {
	return s.status
}

func (s *server) Load() error {
	log.Printf("[%s] Server are loading.", s.config.Name)
	return nil
}

func (s *server) Start() error {
	log.Printf("[%s] Server are starting.", s.config.Name)
	return nil
}

func (s *server) Reload() error {
	log.Printf("[%s] Server are reLoading.", s.config.Name)
	return nil
}

func (s *server) Stop() error {
	log.Printf("[%s] Server are stopping.", s.config.Name)
	return nil
}
