package mFace

type ServerConfig struct {
	Name    string // nickname of the server , default is "Mark"+version
	Network string // type of network about server's listener , default is tcp
	Host    string // default is 0.0.0.0
	Port    string // default is 8888
	MaxConn int64  // default is 0 , 0 means no limit
}

type MConfig interface {
	MServeLoad

	ServerConfig() ServerConfig
	ListenServe() bool
	Address() string
}
