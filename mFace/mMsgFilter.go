package mFace

type MsgFilterType string

var (
	MsgFilterType_Request  MsgFilterType = "Request"
	MsgFilterType_Response MsgFilterType = "Response"
)

type MsgFilterFunc func(MMessage) bool

type MMsgFilter interface {
	Type() MsgFilterType
	FilterFunc() MsgFilterFunc
}
