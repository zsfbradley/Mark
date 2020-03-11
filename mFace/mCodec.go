package mFace

import (
	"net"
)

type CodecCreatorFunc func(net.Conn) MCodec

type MCodec interface {
	ReadRequest() ([]byte, error)
	WriteResponse([]byte) error
	Close() error
}
