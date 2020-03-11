package mNet

import (
	"../mFace"
	"../mTool"
	"errors"
	"net"
)

var (
	Data_Not_Completed_Error = errors.New("data not completed , wait for next round")
)

const (
	Head_Length           = 10
	Message_ID_Length     = 10
	Encryption_Length     = 6
	Length_Of_data_Length = 4

	Head = "_MARK_6_0_"
)

// 10 | 10 | 6 | 4 | n
// 1.10 byte of head
// 2.10 byte of message id
// 3.6 byte of encryption
// 4.4 byte of length about data
func defaultCodec(conn net.Conn) mFace.MCodec {
	cc := &codec{
		conn:   conn,
		buffer: make([]byte, 1024),
	}

	return cc
}

type codec struct {
	conn   net.Conn
	buffer []byte
}

func (cc *codec) ReadRequest() ([]byte, error) {
	if err := cc.acceptData(); err != nil {
		return nil, err
	}

	index, find := cc.readHead()
	if !find {
		return nil, Data_Not_Completed_Error
	}
	lengthOfData, get := cc.readLengthOfData(index)
	if !get {
		return nil, Data_Not_Completed_Error
	}
	data, get := cc.readCompletedData(index, lengthOfData)
	if !get {
		return nil, Data_Not_Completed_Error
	}
	return data, nil
}

func (cc *codec) WriteResponse(responseBytes []byte) error {
	_, err := cc.conn.Write(responseBytes)
	if err != nil {
		return err
	}

	return nil
}

func (cc *codec) Close() error {
	cc.buffer = make([]byte,0)
	cc.buffer = nil
	return cc.conn.Close()
}

// ----------------------------------------------------------- private methods

func (cc *codec) acceptData() error {
	buf := make([]byte, 512)
	cnt, err := cc.conn.Read(buf)
	if err != nil {
		return err
	}
	cc.buffer = append(cc.buffer, buf[:cnt]...)
	return nil
}

func (cc *codec) readHead() (int, bool) {
	length := len(cc.buffer)

	find := false
	index := 0
	for ; index < length; index++ {
		if index+Head_Length > length {
			break
		}

		if string(cc.buffer[index:index+Head_Length]) == Head {
			find = true
			break
		}
	}
	return index, find
}

func (cc *codec) readLengthOfData(offset int) (int, bool) {
	if len(cc.buffer) < offset+Head_Length+Message_ID_Length+Encryption_Length+Length_Of_data_Length {
		return 0, false
	}

	start := offset + Head_Length + Message_ID_Length + Encryption_Length
	end := offset + Head_Length + Message_ID_Length + Encryption_Length + Length_Of_data_Length

	lengthBytes := cc.buffer[start:end]
	length := mTool.ByteToInt(lengthBytes)
	return int(length), true
}

func (cc *codec) readCompletedData(offset, length int) ([]byte, bool) {
	totalLength := Head_Length + Message_ID_Length + Encryption_Length + Length_Of_data_Length + length
	if len(cc.buffer) < (offset + totalLength) {
		return nil, false
	}

	dataBytes := cc.buffer[offset : offset+totalLength]
	return dataBytes, true
}
