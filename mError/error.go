package mError

import (
	"encoding/json"
	"fmt"
)

type Error struct {
	Code int
	Msg  string
}

func (e *Error) IsNil() bool {
	panic("implement me")
}

func (e *Error) Error() string {
	return fmt.Sprintf("MError : %d - %s", e.Code, e.Msg)
}

func (e *Error) TCError() []byte {
	bytes, _ := json.Marshal(e)
	return bytes
}

// 适配 err != nil 的写法，通过 err.NotNil() 判断是否返回错误 ， true说明err非nil，false说明err为nil
func (e *Error) NotNil() bool {
	if e == nil || e.Code == 0 {
		return false
	}

	return true
}
