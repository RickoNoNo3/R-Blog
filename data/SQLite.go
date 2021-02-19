package data

import (
	"github.com/jmoiron/sqlx"

	_ "github.com/mattn/go-sqlite3" // sqlite3
)

var sqldb, sqlerr = sqlx.Open("sqlite3", "./blog.db")

func init() {
	if _, err := sqldb.Exec("PRAGMA recursive_triggers = TRUE"); err != nil {
		panic(err)
	}
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
	err = f(tx)
	if err != nil {
		return
	}
	_ = tx.Commit()
	return
}

// CloseTx close a transaction
func CloseTx(tx *sqlx.Tx) error {
	return tx.Rollback()
}

// CloseDB close the db conn
func CloseDB() {
	if err := sqldb.Ping(); err != nil {
		_ = sqldb.Close()
	}
	return
}
