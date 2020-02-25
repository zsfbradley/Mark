package mMysql

import (
	"markV5/mError"
	"markV5/mFace"
	"sync"
)

var (
	defaultManager = &MysqlManager{
		mysqls:     make(map[string]*mysqlDB),
		mysqlsLock: sync.RWMutex{},
	}
)

func DefaultMysqlManager() *MysqlManager {
	return defaultManager
}

type MysqlManager struct {
	mysqls     map[string]*mysqlDB
	mysqlsLock sync.RWMutex
}

func (mm *MysqlManager) NewMysqlDBConnect(config MysqlConfig) mFace.MError {
	if config.NickName == "" || config.User == "" || config.Port == "" || config.DBName == "" {
		return mError.ErrorParam
	}
	mdb := &mysqlDB{
		config: config,
		db:     nil,
	}

	if err := mdb.Load(); err.NotNil() {
		return err
	}

	if err := mm.register(mdb.config.NickName, mdb); err.NotNil() {
		return err
	}

	return mError.Nil
}

func (mm *MysqlManager) register(nickName string, mdb *mysqlDB) mFace.MError {
	if nickName == "" || mdb == nil {
		return mError.ErrorParam
	}

	mm.mysqlsLock.Lock()
	defer mm.mysqlsLock.Unlock()

	if _, exist := mm.mysqls[nickName]; exist {
		return mError.ParamExistInMap
	}

	mm.mysqls[nickName] = mdb

	return mError.Nil
}

func (mm *MysqlManager) unRegister(nickName string) mFace.MError {
	if nickName == "" {
		return mError.ErrorParam
	}

	mm.mysqlsLock.Lock()
	defer mm.mysqlsLock.Unlock()

	delete(mm.mysqls, nickName)

	return mError.Nil
}

func (mm *MysqlManager) Get(nickName string) (*mysqlDB, mFace.MError) {
	if nickName == "" {
		return nil, mError.ErrorParam
	}

	mm.mysqlsLock.RLock()
	defer mm.mysqlsLock.RUnlock()

	db, exist := mm.mysqls[nickName]
	if !exist {
		return nil, mError.ParamUnExistInMap
	}

	return db, mError.Nil
}
