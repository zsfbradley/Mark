package mNet

import (
	"../mConst"
	"../mFace"
	"../mTool"
	"errors"
	"flag"
	"log"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

var graceful = flag.Bool("graceful", false, "The listener of file descriptor")

func init() {
	flag.Parse()
	log.Printf("This program is taking [%d] pid", os.Getpid())
	log.Printf("The args : %v", os.Args)
}

func NewServer(sConfig ServerConfig) (mFace.MServer, error) {
	sConfig.parse()

	s := &server{
		config:          sConfig,
		status:          mConst.MServe_Status_UnStart,
		listener:        nil,
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
	config          ServerConfig
	status          mConst.MServe_Status
	listener        net.Listener
}

func (s *server) Status() mConst.MServe_Status {
	return s.status
}

func (s *server) Load() error {
	log.Printf("[%s] Server are loading.", s.config.Name)

	// TODO:connManager.Load(),MessageManager.Load(),RouteManager.Load()

	s.status = mConst.MServe_Status_Load

	if err := s.loadListener(); err != nil {
		return err
	}

	log.Printf("[%s] Server finish load.", s.config.Name)
	return nil
}

func (s *server) Start() error {
	log.Printf("[%s] Server are starting.", s.config.Name)

	// TODO:connManager.Start(),MessageManager.Start(),RouteManager.Start()

	s.status = mConst.MServe_Status_Start

	go s.acceptConnect()

	log.Printf("[%s] Server are started.", s.config.Name)

	return s.monitorSignal()
}

func (s *server) Reload() error {
	log.Printf("[%s] Server are reLoading.", s.config.Name)
	return nil
}

func (s *server) Stop() error {
	log.Printf("[%s] Server are stopping.", s.config.Name)

	s.status = mConst.MServe_Status_StartEnding

	_ = s.listener.Close()

	// TODO:connManager.StartEnding(),MessageManager.StartEnding(),RouteManager.StartEnding()
	// TODO:connManager.OfficialEnding(),MessageManager.OfficialEnding(),RouteManager.OfficialEnding()

	s.status = mConst.MServe_Status_OfficialEnding
	s.status = mConst.MServe_Status_Stoped

	log.Printf("[%s] Server are stoped.", s.config.Name)

	return nil
}

// --private functions

func (s *server) loadListener() error {
	if mTool.IsStringEmpty(s.config.Port) {
		return nil // server not gonna listen without port,but still monitor MQ
	}

	var listener net.Listener
	var err error

	if *graceful {
		log.Printf("[%s] Server gonna using file descriptor 3.", s.config.Name)
		f := os.NewFile(3, "")
		listener, err = net.FileListener(f)
	} else {
		log.Printf("[%s] Server gonna using new file descriptor.", s.config.Name)
		addr := net.JoinHostPort(s.config.Host, s.config.Port)
		listener, err = net.Listen(s.config.Network, addr)
	}

	if err != nil {
		return err
	}

	s.listener = listener

	return nil
}

func (s *server) acceptConnect() {
	if s.listener == nil {
		return
	}

	log.Printf("[%s] Server accept connect on %s", s.config.Name, s.listener.Addr())

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.status >= mConst.MServe_Status_StartEnding {
				break
			}
			log.Printf("[%s] Listener accept connect error : %v", s.config.Name, err)
			continue
		}
		go handleConn(conn)
	}
	log.Printf("[%s] Server are out of accept connect.", s.config.Name)
}

func handleConn(conn net.Conn) {

}

func (s *server) monitorSignal() error {
	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, syscall.SIGHUP, syscall.SIGINT,
		syscall.SIGTERM, syscall.SIGUSR1, syscall.SIGUSR2)
	signalInfo := <-signalChan

	log.Printf("[%s] Server accept signal info : %s", s.config.Name, signalInfo)

	//MonitorSignal:
	switch signalInfo {
	case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM:
		log.Printf("[%s] Server gonna stop", s.config.Name)
		if err := s.Stop(); err != nil {
			return err
		}
		signal.Stop(signalChan)
	//case syscall.SIGUSR1:
	//	log.Printf("[%s] Server gonna reload", s.config.Name)
	//	if err := s.Reload(); err != nil {
	//		return err
	//	}
	//	goto MonitorSignal
	case syscall.SIGUSR2:
		log.Printf("[%s] Server gonna restart", s.config.Name)
		signal.Stop(signalChan)
		if err := s.restart(); err != nil {
			return err
		}
	}

	return nil
}

func (s *server) restart() error {
	tcpListener, ok := s.listener.(*net.TCPListener)
	if !ok {
		return errors.New("listener is not tcp listener")
	}

	f, err := tcpListener.File()
	if err != nil {
		return err
	}

	args := []string{"-graceful"}
	cmd := exec.Command(os.Args[0], args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.ExtraFiles = []*os.File{f}

	return cmd.Start()
}
