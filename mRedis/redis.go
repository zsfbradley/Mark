package mRedis

import (
	"github.com/garyburd/redigo/redis"
	"markV5/mError"
	"markV5/mFace"
	"time"
)

type MRedis struct {
	config MRedisConfig
	pool   *redis.Pool
}

func (r *MRedis) load() mFace.MError {
	r.pool = &redis.Pool{
		Dial: func() (conn redis.Conn, err error) {
			return redis.Dial(
				r.config.Network,
				r.config.address(),
				redis.DialReadTimeout(time.Second*time.Duration(r.config.ReadTimeout)),
				redis.DialWriteTimeout(time.Second*time.Duration(r.config.WriteTimeout)),
				redis.DialConnectTimeout(time.Second*time.Duration(r.config.ConnectTimeout)),
			)
		},
		MaxIdle:     r.config.PoolMaxIdle,
		MaxActive:   r.config.PoolMaxActive,
		IdleTimeout: time.Second * time.Duration(r.config.PoolIdleTimeout),
		Wait:        r.config.PoolWait,
	}

	return mError.Nil
}

func (r *MRedis) get() (redis.Conn, mFace.MError) {
	conn := r.pool.Get() // 获取连接
	if err := conn.Err(); err != nil {
		return nil, mError.SystemError(err)
	}
	return conn, mError.Nil
}

func (r *MRedis) Exec(cmd string, key interface{}, args ...interface{}) (interface{}, mFace.MError) {
	params := make([]interface{}, 0)
	params = append(params, key)

	if len(args) > 0 {
		for _, param := range args {
			params = append(params, param)
		}
	}

	return r.baseExec(cmd, params...)
}

func (r *MRedis) baseExec(cmd string, args ...interface{}) (interface{}, mFace.MError) {
	conn, err := r.get() // 获取连接
	if err.NotNil() {
		return nil, err
	}
	defer conn.Close()

	result, doErr := conn.Do(cmd, args...)
	if doErr != nil {
		return result, mError.SystemError(doErr)
	}

	return result, mError.Nil
}

// 查看所有键
func (r *MRedis) Keys() (interface{}, mFace.MError) {
	return r.Exec("KEYS", "*")
}

// 查看所有键
func (r *MRedis) DBSize() (interface{}, mFace.MError) {
	return r.baseExec("dbsize")
}

// 查看 key 键是否存在
func (r *MRedis) Exists(key string) (interface{}, mFace.MError) {
	return r.Exec("EXISTS", key)
}

// 删除 keys
func (r *MRedis) DelKeys(keys ...string) (interface{}, mFace.MError) {
	params := make([]interface{}, 0)
	for _, key := range keys {
		params = append(params, key)
	}
	return r.baseExec("DEL", params...)
}

// 给 key 指定过期时间，单位秒
func (r *MRedis) Expire(key string, seconds int) (interface{}, mFace.MError) {
	return r.Exec("EXPIRE", key, seconds)
}

// 给 key 以时间戳指定过期时间，单位秒，seconds 为时间戳
func (r *MRedis) Expireat(key string, seconds int64) (interface{}, mFace.MError) {
	return r.Exec("EXPIREAT", key, seconds)
}

// 给 key 指定过期时间，单位毫秒
func (r *MRedis) Pexpire(key string, milliSeconds int64) (interface{}, mFace.MError) {
	return r.Exec("PEXPIRE", key, milliSeconds)
}

// 给 key 以时间戳指定过期时间，单位毫秒，milliSeconds 为时间戳
func (r *MRedis) Pexpireat(key string, milliSeconds int64) (interface{}, mFace.MError) {
	return r.Exec("EXPIREAT", key, milliSeconds)
}

// 查询 key 的过期时间，单位秒
// 大于等于0，剩余过期时间
// -1 ，键未设置过期时间
// -2 ，键不存在
func (r *MRedis) TTL(key string) (interface{}, mFace.MError) {
	return r.Exec("TTL", key)
}

// 查询 key 的过期时间，单位毫秒
// 大于等于0，剩余过期时间
// -1 ，键未设置过期时间
// -2 ，键不存在
func (r *MRedis) PTTL(key string) (interface{}, mFace.MError) {
	return r.Exec("PTTL", key)
}

// 清除 key 的过期时间
func (r *MRedis) Persist(key string) (interface{}, mFace.MError) {
	return r.Exec("PERSIST", key)
}

// 查询 key 的类型
func (r *MRedis) Type(key string) (interface{}, mFace.MError) {
	return r.Exec("TYPE", key)
}

// key 重命名
func (r *MRedis) Rename(key, newKey string) (interface{}, mFace.MError) {
	return r.Exec("RENAME", key, newKey)
}

// key 重命名，确保 newKey 不存在才覆盖 key
func (r *MRedis) RenameNX(key, newKey string) (interface{}, mFace.MError) {
	return r.Exec("RENAMENX", key, newKey)
}

// 随机返回一个键
func (r *MRedis) RandomKey() (interface{}, mFace.MError) {
	return r.baseExec("RANDOMKEY")
}
