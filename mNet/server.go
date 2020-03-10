package mNet

import (
	"../mConst"
	"../mFace"
	"log"
	"net"
)

func init() {
	log.Println(mConst.Framework_Line)
	log.Printf("%s[%s%s]%s", mConst.Framework_Line_Half, mConst.Framework_Name, mConst.Framework_Version, mConst.Framework_Line_Half)
	log.Println(mConst.Framework_Line)
}

func NewServer() mFace.MServer {
	return NewServerWithConfig(defaultConfig())
}

func NewServerWithConfigPath(filePath string) mFace.MServer {
	return NewServerWithConfig(loadConfigWithFilePath(filePath))
}

func NewServerWithConfig(config mFace.MConfig) mFace.MServer {
	s := &server{
		status:   mConst.MServe_Status_UnStart,
		cm:       newConnManager(),
		mm:       newMessageManager(),
		rm:       newRouteManager(),
		config:   config,
		listener: nil,
	}

	s.cm.BindServer(s)
	s.mm.BindServer(s)
	s.rm.BindServer(s)

	return s
}

type server struct {
	status mConst.MServe_Status
	cm     mFace.MConnManager
	mm     mFace.MMessageManager
	rm     mFace.MRouteManager

	config   mFace.MConfig
	listener net.Listener
}

func (s *server) Status() mConst.MServe_Status {
	return s.status
}

func (s *server) Config() mFace.ServerConfig {
	return s.config.ServerConfig()
}

func (s *server) Load() error {
	s.status = mConst.MServe_Status_Load
	log.Printf("[%s] start load", mConst.Framework_Name+mConst.Framework_Version)

	if err := s.config.Load(); err != nil {
		return err
	}

	if err := s.cm.Load(); err != nil {
		return err
	}

	if err := s.mm.Load(); err != nil {
		return err
	}

	if err := s.rm.Load(); err != nil {
		return err
	}

	// load listener in the end
	if s.config.ListenServe() {
		listener, err := net.Listen(s.config.ServerConfig().Network, s.config.Address())
		if err != nil {
			return err
		}

		s.listener = listener
	}

	log.Printf("[%s] finish load", s.config.ServerConfig().Name)

	return nil
}

func (s *server) Start() error {
	log.Printf("[%s] starting", s.config.ServerConfig().Name)

	if err := s.cm.Start(); err != nil {
		return err
	}

	if err := s.mm.Start(); err != nil {
		return err
	}

	if err := s.rm.Start(); err != nil {
		return err
	}

	s.status = mConst.MServe_Status_Start

	// start accept connection
	if s.listener != nil && s.config.ListenServe() {
		go s.startAcceptConnection()
	}

	log.Printf("[%s] started", s.config.ServerConfig().Name)
	return nil
}

func (s *server) Reload() error {
	s.status = mConst.MServe_Status_Reload
	log.Printf("[%s] start reload", mConst.Framework_Name+mConst.Framework_Version)

	lastListenAddr := s.config.Address()

	if err := s.config.Reload(); err != nil {
		return err
	}

	if err := s.cm.Reload(); err != nil {
		return err
	}

	if err := s.mm.Reload(); err != nil {
		return err
	}

	if err := s.rm.Reload(); err != nil {
		return err
	}

	// reload and start listener in the end
	if s.config.ListenServe() && lastListenAddr != s.config.Address() {
		listener, err := net.Listen(s.config.ServerConfig().Network, s.config.Address())
		if err != nil {
			return err
		}

		// close last listener
		if s.listener != nil {
			if err := s.listener.Close(); err != nil {
				return err
			}
		}

		s.listener = listener

		go s.startAcceptConnection()
	}

	log.Printf("[%s] finish reload", s.config.ServerConfig().Name)

	s.status = mConst.MServe_Status_Start

	return nil
}

func (s *server) Stop() error {
	log.Printf("[%s] start stop", s.config.ServerConfig().Name)

	s.status = mConst.MServe_Status_StartEnding

	if s.listener != nil {
		if err := s.listener.Close(); err != nil {
			return err
		}
	}

	if err := s.cm.StartEnding(); err != nil {
		return err
	}

	if err := s.mm.StartEnding(); err != nil {
		return err
	}

	if err := s.rm.StartEnding(); err != nil {
		return err
	}

	if err := s.rm.OfficialEnding(); err != nil {
		return err
	}

	if err := s.cm.OfficialEnding(); err != nil {
		return err
	}

	if err := s.rm.OfficialEnding(); err != nil {
		return err
	}

	s.status = mConst.MServe_Status_OfficialEnding
	s.status = mConst.MServe_Status_Stoped
	log.Printf("[%s] stoped", s.config.ServerConfig().Name)

	return nil
}

// ----------------------------------------------------------- private methods

func (s *server) startAcceptConnection() {
	log.Printf("[%s] server gonna listen and serve on %s", s.config.ServerConfig().Name, s.config.Address())
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Printf("[%s] server listener on %saccept error %s", s.config.ServerConfig().Name, s.config.Address(), err)
			if s.status >= mConst.MServe_Status_StartEnding {
				break
			} else {
				continue
			}
		}
		go handleConn(conn)
	}

	log.Printf("[%s] server stop listen and serve on %s", s.config.ServerConfig().Name, s.config.Address())
}

func handleConn(conn net.Conn) {

}
