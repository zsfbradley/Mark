package mConst

type MServeStatus uint8

const (
	MServeStatus_UnStart = iota
	MServeStatus_Load
	MServeStatus_Start
	MServeStatus_Reload
	MServeStatus_StartEnding
	MServeStatus_OfficialEnding
	MServeStatus_Stopped
)
