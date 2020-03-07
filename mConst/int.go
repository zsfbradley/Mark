package mConst

type MServe_Status int

const (
	MServe_Status_UnStart        MServe_Status = iota
	MServe_Status_Load           MServe_Status = iota << 1
	MServe_Status_Start          MServe_Status = iota << 1
	MServe_Status_Reload         MServe_Status = iota << 1
	MServe_Status_StartEnding    MServe_Status = iota << 1
	MServe_Status_OfficialEnding MServe_Status = iota << 1
	MServe_Status_Stoped         MServe_Status = iota << 1
)
