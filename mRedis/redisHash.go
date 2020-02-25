package mRedis

import "markV5/mFace"

/*
经典使用场景有:
1.用户信息缓存:将用户信息的属性及值一一对应，放入一个键中，通过固定前、后缀+唯一识别标志来作为键
*/
// 设置 key 的一对键值对
func (r *MRedis) HashSet(key string, field, value interface{}) (interface{}, mFace.MError) {
	return r.Exec("HSET", key, field, value)
}

// 设置 key 的N对键值对
func (r *MRedis) HashMSet(key string, maps map[interface{}]interface{}) (interface{}, mFace.MError) {
	params := make([]interface{}, 0)
	for sKey, value := range maps {
		params = append(params, sKey)
		params = append(params, value)
	}
	return r.Exec("HMSET", key, params...)
}

// 获取 key 的指定键的值
func (r *MRedis) HashGet(key string, field interface{}) (interface{}, mFace.MError) {
	return r.Exec("HGET", key, field)
}

// 获取 key 的N对键值对
func (r *MRedis) HashMGet(key string, fields ...interface{}) (interface{}, mFace.MError) {
	params := make([]interface{}, 0)
	for _, sKey := range fields {
		params = append(params, sKey)
	}
	return r.Exec("HMGET", key, params...)
}

// 删除 key 的N对键值对
func (r *MRedis) HashDelete(key string, fields ...interface{}) (interface{}, mFace.MError) {
	params := make([]interface{}, 0)
	params = append(params, key)
	params = append(params, fields...)
	return r.Exec("HDEL", key, fields...)
}

// 统计 key 的键值对数
func (r *MRedis) HashLen(key string) (interface{}, mFace.MError) {
	return r.Exec("HLEN", key)
}

// 判断 key 的键 field 是否存在 , 不存在返回 0 ， 存在返回 1
func (r *MRedis) HashExists(key string, field interface{}) (interface{}, mFace.MError) {
	return r.Exec("HEXISTS", key, field)
}

// 获取 key 的所有键
func (r *MRedis) HashKeys(key string) (interface{}, mFace.MError) {
	return r.Exec("HKEYS", key)
}

// 获取 key 的所有的键的值
func (r *MRedis) HashValues(key string) (interface{}, mFace.MError) {
	return r.Exec("HVALS", key)
}

// 获取 key 的所有的键与值 , 返回的数据按 []string{键1，值1，键2，值2...} 的顺序
// 如果 key 的键值对较多，使用 hscan 命令
func (r *MRedis) HashGetAll(key string) (interface{}, mFace.MError) {
	return r.Exec("HGETALL", key)
}

// hash 指定 key 的 field 字段以 increment 自增一次
func (r *MRedis) HashIncrBy(key string, field interface{}, increment int) (interface{}, mFace.MError) {
	return r.Exec("HINCRBY", key, field, increment)
}

// hash 指定 key 的 field 字段以 increment 浮点数自增一次
// 使用 HashIncrByFloat 自增过后的字段再使用 HashIncrBy 自增 integer 会发生错误
func (r *MRedis) HashIncrByFloat(key string, field interface{}, increment float64) (interface{}, mFace.MError) {
	return r.Exec("HINCRBYFLOAT", key, field, increment)
}

// 获取 key 的键 field 的长度
func (r *MRedis) HashStrLen(key string, field interface{}) (interface{}, mFace.MError) {
	return r.Exec("HSTRLEN", key, field)
}
