package mNet

import "markV5/mFace"

func newMsgFilter(msgType mFace.MsgFilterType, msgFilterFunc mFace.MsgFilterFunc) mFace.MMsgFilter {
	mf := &MsgFilter{
		msgType:       msgType,
		msgFilterFunc: msgFilterFunc,
	}
	return mf
}

type MsgFilter struct {
	msgType       mFace.MsgFilterType
	msgFilterFunc mFace.MsgFilterFunc
}

func (mf *MsgFilter) Type() mFace.MsgFilterType {
	return mf.msgType
}

func (mf *MsgFilter) FilterFunc() mFace.MsgFilterFunc {
	return mf.msgFilterFunc
}
