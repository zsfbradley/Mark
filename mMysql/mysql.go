package mMysql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"markV5/mError"
	"markV5/mFace"
)

const (
	Driver_Mysql = "mysql"
)

type mysqlDB struct {
	config MysqlConfig
	db     *sql.DB
}

func (mdb *mysqlDB) Load() mFace.MError {
	db, err := sql.Open(Driver_Mysql, mdb.config.dsn())
	if err != nil {
		return mError.SystemError(err)
	}

	if err := db.Ping(); err != nil {
		return mError.SystemError(err)
	}

	mdb.db = db

	return mError.Nil
}

func (mdb *mysqlDB) Close() mFace.MError {
	if mdb.db != nil {
		err := mdb.db.Close()
		if err != nil {
			return mError.SystemError(err)
		}

		// 主动关闭，从管理类中删除
		if err := DefaultMysqlManager().unRegister(mdb.config.NickName); err.NotNil() {
			return err
		}
	}

	return mError.Nil
}

func (mdb *mysqlDB) DB() *sql.DB {
	return mdb.db
}

func (mdb *mysqlDB) CreateDatabase(name string) mFace.MError {
	_, err := mdb.db.Exec(fmt.Sprintf("CREATE DATABASE %s;", name))
	if err != nil {
		return mError.SystemError(err)
	}
	return mError.Nil
}

func (mdb *mysqlDB) DropDatabase(name string) mFace.MError {
	_, err := mdb.db.Exec(fmt.Sprintf("DROP DATABASE %s;", name))
	if err != nil {
		return mError.SystemError(err)
	}
	return mError.Nil
}
