package data

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3" // sqlite3
	"os"
	"rickonono3/r-blog/logger"

	"rickonono3/r-blog/objects"
)

var sqldb *sqlx.DB
var sqlerr error

func OpenDB(dbLoc string) {
	file := objects.CWD + dbLoc
	_, err := os.Stat(file)
	if err != nil {
		logger.L.Panic("[Database]", "无法加载数据库!")
		panic("无法加载数据库!")
	}
	sqldb, sqlerr = sqlx.Open(
		"sqlite3",
		file,
	)
	sqldb.MustExec("PRAGMA journal_mode = DELETE")
	sqldb.MustExec("PRAGMA recursive_triggers = TRUE")
}

// NewTx get a new sqlite transaction
// use defer to CloseTx(CloseDB as well)
func NewTx() (tx *sqlx.Tx, err error) {
	if sqlerr != nil {
		err = sqlerr
		return
	}
	tx, err = sqldb.Beginx()
	return
}

// DoTx is an easy way to begin and do something in transactions
func DoTx(f func(*sqlx.Tx) (err error)) (err error) {
	tx, err := NewTx()
	if err != nil {
		return
	}
	defer CloseTx(tx)
	logger.L.Info("[Database]", "数据库事务-开启")
	err = f(tx)
	if err != nil {
		return
	}
	err = tx.Commit()
	if err == nil {
		logger.L.Info("[Database]", "数据库事务-提交")
	}
	return
}

// CloseTx close a transaction
func CloseTx(tx *sqlx.Tx) (err error) {
	err = tx.Rollback()
	tx.Rollback()
	if err == nil {
		logger.L.Info("[Database]", "数据库事务-回滚")
	}
	return
}

// CloseDB close the db conn
func CloseDB() {
	if err := sqldb.Ping(); err != nil {
		_ = sqldb.Close()
	}
	return
}
