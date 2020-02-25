package mError

import (
	"fmt"
	"markV5/mFace"
)

func SystemError(err error) mFace.MError {
	return &Error{
		Code: 5000,
		Msg:  fmt.Sprintf("Serve get an Error : %s.Please reply to customer service.", err.Error()),
	}
}

var (
	Nil               = &Error{Code: 0, Msg: ""}
	ErrorParam        = &Error{Code: 5001, Msg: "param error"}
	NilParam          = &Error{Code: 5002, Msg: "param is nil"}
	RouteExist        = &Error{Code: 5003, Msg: "route id exist"}
	RouteUnExist      = &Error{Code: 5004, Msg: "route id not exist"}
	ChannelUnOpen     = &Error{Code: 5004, Msg: "channel status not start or reload"}
	ConnIDUnExist     = &Error{Code: 5005, Msg: "Conn id not exist"}
	ParamExistInMap   = &Error{Code: 5006, Msg: "param exist in map"}
	ParamUnExistInMap = &Error{Code: 5007, Msg: "param unExist in map"}

	InsertDataMysql = &Error{Code: 5100, Msg: "insert data into database error"}
	InsertDataRedis = &Error{Code: 5200, Msg: "insert data into redis error"}
)
