package mRedis

import "markV5/mFace"

/*
经典使用场景有:
1.计数:视频观看数，点赞数等自增计数工具
*/

// 一旦对某integer键进行了浮点数自增后，再使用Incr、IncrBy会报错，因为此时该键不再是integer类型了

// 自增1，如果key不存在，默认为0，自增1返回
func (r *MRedis) Incr(key string) (interface{}, mFace.MError) {
	return r.Exec("INCR", key)
}

// 自减1，如果key不存在，默认为0，自减1返回
func (r *MRedis) Decr(key string) (interface{}, mFace.MError) {
	return r.Exec("DECR", key)
}

// 以指定的increment参数自增一次
func (r *MRedis) IncrBy(key string, increment int) (interface{}, mFace.MError) {
	return r.Exec("INCRBY", key, increment)
}

// 以指定的decrement参数自减一次
func (r *MRedis) DecrBy(key string, decrement int) (interface{}, mFace.MError) {
	return r.Exec("DECRBY", key, decrement)
}

// 以指定的increment浮点参数自增一次
func (r *MRedis) IncrByFloat(key string, increment float64) (interface{}, mFace.MError) {
	return r.Exec("INCRBYFLOAT", key, increment)
}
