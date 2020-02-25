package mNet

import (
	"markV5/mConst"
	"markV5/mError"
	"markV5/mFace"
)

func newChannel() mFace.MChannel {
	c := &channel{
		status:   mConst.MServeStatus_UnStart,
		size:     0,
		realChan: nil,
	}
	return c
}

type channel struct {
	status   mConst.MServeStatus
	size     uint64
	realChan *chan interface{}
}

func (c *channel) Load() mFace.MError {
	c.status = mConst.MServeStatus_Load
	realChan := make(chan interface{}, c.size)
	c.realChan = &realChan
	return nil
}

func (c *channel) Reload() mFace.MError {
	c.status = mConst.MServeStatus_Reload
	newChan := make(chan interface{}, c.size)
	oldChan := c.realChan
	c.realChan = &newChan
	if len(*oldChan) > 0 {
		go func() {
			for {
				if len(*oldChan) == 0 {
					break
				}
				data := <-*oldChan
				*(c.realChan) <- data
			}
		}()
	}
	return nil
}

func (c *channel) Start() mFace.MError {
	c.status = mConst.MServeStatus_Start
	return nil
}

func (c *channel) Status() mConst.MServeStatus {
	return c.status
}

func (c *channel) StartEnding() mFace.MError {
	c.status = mConst.MServeStatus_StartEnding
	for {
		if len(*c.realChan) == 0 {
			break
		}
	}
	return nil
}

func (c *channel) OfficialEnding() mFace.MError {
	c.status = mConst.MServeStatus_OfficialEnding
	close(*c.realChan)
	return nil
}

func (c *channel) SetSize(size uint64) {
	c.size = size
}

func (c *channel) In(data interface{}) mFace.MError {
	if c.status < mConst.MServeStatus_Start || c.status > mConst.MServeStatus_Reload { // start 、 reload 模式下可以继续存入数据
		return mError.ChannelUnOpen
	}
	*c.realChan <- data

	return mError.Nil
}

func (c *channel) Out() (interface{}, bool) {
	if c.status < mConst.MServeStatus_Start || c.status >= mConst.MServeStatus_OfficialEnding { // 非 unStart 、Load 、 officialEnding 模式下可以取得数据
		return nil, false
	}

	data, ok := <-(*c.realChan)
	if !ok {
		return nil, false
	}
	return data, true
}

func (c *channel) Close() mFace.MError {
	c.status = mConst.MServeStatus_Stopped // 状态强制置为已结束
	close(*c.realChan)
	return nil
}
