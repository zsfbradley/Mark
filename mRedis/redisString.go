package mRedis

import (
	"markV5/mError"
	"markV5/mFace"
)

/*
经典使用场景有:
1.缓存功能：redis作为缓存层，mysql作为存储层，将登录请求的用户信息获取验证通过后，存入redis中，设置过期时间
2.共享session：当多个web服务公用一套session逻辑时，可以将用户的session缓存到redis中，由多个web服务共享
*/

type Expire_Type string
type Action_Type string

const (
	Expire_Nil           Expire_Type = ""
	Expire_Seconds       Expire_Type = "ex"
	Expire_MillisSeconds Expire_Type = "px"
	Action_Nil           Action_Type = ""
	Action_Add           Action_Type = "nx"
	Action_Update        Action_Type = "xx"
)

// 设置键（必选项）
func (r *MRedis) SetWithOption(key string, value interface{}, expireType Expire_Type, expireValue int, actionType Action_Type) (interface{}, mFace.MError) {
	params := make([]interface{}, 0)
	if value == nil {
		return nil, mError.ErrorParam
	}
	params = append(params, value)
	if expireType != Expire_Nil {
		params = append(params, expireType)
		params = append(params, expireValue)
	}
	if actionType != Action_Nil {
		params = append(params, actionType)
	}
	return r.Exec("SET", key, params...)
}

// 批量设置键
func (r *MRedis) MSet(maps map[string]interface{}) (interface{}, mFace.MError) {
	params := make([]interface{}, 0)
	for key, value := range maps {
		params = append(params, key)
		params = append(params, value)
	}

	return r.baseExec("MSET", params...)
}

// 获取值
func (r *MRedis) Get(key string, args ...interface{}) (interface{}, mFace.MError) {
	return r.Exec("GET", key, args...)
}

// 批量获取值
func (r *MRedis) MGet(keys ...string) (interface{}, mFace.MError) {
	params := make([]interface{}, 0)
	for _, key := range keys {
		params = append(params, key)
	}
	return r.baseExec("MGET", params...)
}

// 向 key 的值追加字符串，如果 key 不存在则创建
func (r *MRedis) Append(key string, value string) (interface{}, mFace.MError) {
	return r.Exec("APPEND", key, value)
}

// 查询 key 的值长度
func (r *MRedis) StrLen(key string) (interface{}, mFace.MError) {
	return r.Exec("STRLEN", key)
}

// 设置 key ， 返回原来的 value
func (r *MRedis) GetSet(key string, value interface{}) (interface{}, mFace.MError) {
	return r.Exec("GETSET", key, value)
}

// 设置 key 指定从 location 位置起的值为 value
func (r *MRedis) SetRange(key string, location int, value interface{}) (interface{}, mFace.MError) {
	return r.Exec("SETRANGE", key, location, value)
}

// 获取 key 指定从 start 位置至 end 位置的值（包括end位），若长度超出则至结尾为止
func (r *MRedis) GetRange(key string, start, end int) (interface{}, mFace.MError) {
	return r.Exec("GETRANGE", key, start, end)
}
