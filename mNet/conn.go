package mNet

import (
	"io"
	"log"
	"markV5/mConst"
	"markV5/mError"
	"markV5/mFace"
	"net"
)

var (
	stages = []uint64{10, 4} // 自定义数据协议段落 10-路由ID 4-实际数据长度记录
)

const (
	symbol = "HEADERSYMBOL" // 自定义数据协议标志
)

func newConn(id string, oConn *net.Conn) mFace.MConn {
	c := &conn{
		id:              id,
		conn:            oConn,
		cm:              nil,
		status:          mConst.MServeStatus_UnStart,
		dp:              nil,
		startEndingChan: make(chan bool, 0),
	}
	return c
}

type conn struct {
	id              string
	conn            *net.Conn
	status          mConst.MServeStatus
	dp              mFace.MDataProtocol
	cm              mFace.MConnManager
	startEndingChan chan bool
}

func (c *conn) SuperiorManager(cm mFace.MConnManager) {
	c.cm = cm
}

func (c *conn) Load() mFace.MError {
	c.status = mConst.MServeStatus_Load

	dp, err := newDataProtocol(symbol, stages...)
	if err.NotNil() {
		return err
	}
	c.dp = dp

	return mError.Nil
}

func (c *conn) Reload() mFace.MError {
	c.status = mConst.MServeStatus_Reload
	return mError.Nil
}

func (c *conn) Start() mFace.MError {
	log.Printf("[%s] start comunicate", c.id)
	c.status = mConst.MServeStatus_Start

	go c.acceptCompletedRequest()
	go c.acceptRequestData()

	return mError.Nil
}

func (c *conn) Status() mConst.MServeStatus {
	return c.status
}

func (c *conn) StartEnding() mFace.MError {
	log.Printf("[%s] start ending comunicate", c.id)
	c.status = mConst.MServeStatus_StartEnding
	// 关闭接收
	c.startEndingChan <- true
	return mError.Nil
}

func (c *conn) OfficialEnding() mFace.MError {
	log.Printf("[%s] official comunicate", c.id)
	c.status = mConst.MServeStatus_OfficialEnding
	// 关闭回复
	if err := (*(c.conn)).Close(); err != nil {
		return mError.SystemError(err)
	}
	log.Printf("[%s] official stopped", c.id)
	return mError.Nil
}

func (c *conn) ConnID() string {
	return c.id
}

func (c *conn) Response(rsp mFace.MMessage) {
	if c.status >= mConst.MServeStatus_OfficialEnding {
		return
	}

	rConn := *(c.conn)
	if _, err := rConn.Write(rsp.Marshal()); err != nil {
		log.Println(err.Error())
	}
}

// - private methods

func (c *conn) acceptRequestData() {
	rConn := *(c.conn)
	for {
		buf := make([]byte, 512)
		cnt, err := rConn.Read(buf)
		if c.status >= mConst.MServeStatus_StartEnding {
			break
		}
		if err != nil {
			if err == io.EOF { // 用户断开连接
				c.breakConnection()
				break
			}
			continue
		}
		c.dp.Unmarshal(buf[:cnt])
	}
	log.Printf("[%s] Conn End Of Accept Request Data", c.id)
}

func (c *conn) acceptCompletedRequest() {
	for {
		breakMonitor := false
		select {
		case datas := <-c.dp.CompletedDataChannel():
			routeID := ""
			if len(datas) >= 2 {
				routeID = string(datas[1])
			}
			data := datas[len(datas)-1]
			request := newMessage(routeID, c.id, data)
			if err := c.cm.DataInChannel().In(request); err.NotNil() {
				log.Println(err.Error())
			}
		case <-c.startEndingChan:
			breakMonitor = true
			break
		}
		if breakMonitor {
			break
		}
	}
	log.Printf("[%s] Conn End Of Accept Completed Request", c.id)
}

func (c *conn) breakConnection() {
	c.startEndingChan <- true
	if err := c.cm.ConnOut(c.id); err.NotNil() {
		log.Println(err.Error())
	}
	log.Printf("[%s] Break Connection", c.id)
}
