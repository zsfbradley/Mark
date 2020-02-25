package mRedis

import (
	"markV5/mError"
	"markV5/mFace"
	"sync"
)

var (
	defaultRedisManager = &RedisManager{
		redisMap:  make(map[string]*MRedis),
		redisLock: sync.RWMutex{},
	}
)

func DefaultRedisManager() *RedisManager {
	return defaultRedisManager
}

type RedisManager struct {
	redisMap  map[string]*MRedis
	redisLock sync.RWMutex
}

func (rm *RedisManager) NewRedisConnect(config MRedisConfig) mFace.MError {
	if config.Nickname == "" || config.Network == "" || config.Host == "" || config.Port == "" {
		return mError.ErrorParam
	}

	r := &MRedis{
		config: config,
		pool:   nil,
	}

	if err := r.load(); err.NotNil() {
		return err
	}

	if err := defaultRedisManager.register(r.config.Nickname, r); err.NotNil() {
		return err
	}

	return mError.Nil
}

func (rm *RedisManager) register(nickName string, r *MRedis) mFace.MError {
	if nickName == "" || r == nil {
		return mError.ErrorParam
	}

	rm.redisLock.Lock()
	defer rm.redisLock.Unlock()

	_, exist := rm.redisMap[nickName]
	if exist {
		return mError.ParamExistInMap
	}

	rm.redisMap[nickName] = r

	return mError.Nil
}

func (rm *RedisManager) Get(nickName string) (*MRedis, mFace.MError) {
	if nickName == "" {
		return nil, mError.ErrorParam
	}

	rm.redisLock.RLock()
	defer rm.redisLock.RUnlock()

	r, exist := rm.redisMap[nickName]
	if !exist {
		return nil, mError.ParamUnExistInMap
	}

	return r, mError.Nil
}
