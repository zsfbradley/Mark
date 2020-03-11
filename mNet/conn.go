package mNet

import (
	"../mConst"
	"../mFace"
	"fmt"
	"io"
	"log"
)

func newConn(id string, codec mFace.MCodec) mFace.MConn {
	c := &conn{
		status:  mConst.MServe_Status_UnStart,
		id:      id,
		codec:   codec,
		manager: nil,
	}
	return c
}

type conn struct {
	status  mConst.MServe_Status
	id      string
	codec   mFace.MCodec
	manager mFace.MConnManager
}

func (c *conn) BindManager(manager mFace.MConnManager) {
	c.manager = manager
}

func (c *conn) Status() mConst.MServe_Status {
	return c.status
}

func (c *conn) Load() error {
	c.status = mConst.MServe_Status_Load
	log.Printf("[%s] conn start load", c.id)

	log.Printf("[%s] conn finish load", c.id)

	return nil
}

func (c *conn) Start() error {
	c.status = mConst.MServe_Status_Start
	log.Printf("[%s] conn starting", c.id)

	go c.startAcceptRequest()

	log.Printf("[%s] conn started", c.id)

	return nil
}

func (c *conn) Reload() error {
	c.status = mConst.MServe_Status_Reload
	log.Printf("[%s] conn start reload", c.id)

	log.Printf("[%s] conn finish reload", c.id)
	c.status = mConst.MServe_Status_Start

	return nil
}

func (c *conn) StartEnding() error {
	c.status = mConst.MServe_Status_StartEnding
	log.Printf("[%s] conn start ending", c.id)
	return nil
}

func (c *conn) OfficialEnding() error {
	c.status = mConst.MServe_Status_OfficialEnding
	log.Printf("[%s] conn official ending", c.id)

	c.status = mConst.MServe_Status_Stoped
	log.Printf("[%s] conn stoped", c.id)

	return nil
}

func (c *conn) ID() string {
	return c.id
}

func (c *conn) ReplyResponse(responseBytes []byte) error {
	return c.codec.WriteResponse(responseBytes)
}

// ----------------------------------------------------------- private methods

func (c *conn) startAcceptRequest() {
	for {
		requestBytes, err := c.codec.ReadRequest()
		if err != nil {
			if err == io.EOF {
				c.breakConnect()
				break
			}
			if c.status >= mConst.MServe_Status_StartEnding {
				break
			}
			continue
		}
		err = c.ReplyResponse([]byte(fmt.Sprintf("[%s] get you message : %s", c.id, string(requestBytes))))
		if err != nil {
			log.Println(err)
		}
		// todo : send requestBytes to connManager.InputChan
	}
}

func (c *conn) breakConnect() {
	log.Printf("[%s] conn break connect", c.id)
	log.Printf("[%s] conn close codec : %v", c.id, c.codec.Close())
	c.manager.DeleteConn(c.id)
}
