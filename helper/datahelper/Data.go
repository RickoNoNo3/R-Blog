package datahelper

import (
	"database/sql"

	"github.com/jmoiron/sqlx"

	"rickonono3/r-blog/helper/typehelper"
	"rickonono3/r-blog/mytype"
)

// 将 res 转为 int 型的 id
func GetResultId(res *sql.Result) (id int) {
	id64, _ := (*res).LastInsertId()
	id = typehelper.Int64ToInt(id64)
	return
}

// 判断某个实体是否存在
func IsExists(tx *sqlx.Tx, entity mytype.Entity) bool {
	if entity.Type == 0 && entity.Id == 0 {
		return true
	}
	var (
		err          error
		count        = 0
		layerExists  = false
		entityExists = false
	)
	// layerExists
	if err = tx.QueryRowx(
		"select count(*) from layer where type=? and id=?",
		entity.Type,
		entity.Id,
	).Scan(&count); err != nil {
		return false
	}
	if count == 1 {
		layerExists = true
	}
	// entityExists
	count = 0
	query := ""
	switch entity.Type {
	case 0:
		query = "select count(*) from dir where id=?"
	case 1:
		query = "select count(*) from article where id=?"
	case 2:
		query = "select count(*) from file where id=?"
	}
	if err = tx.QueryRowx(
		query,
		entity.Id,
	).Scan(&count); err != nil {
		return false
	}
	if count == 1 {
		entityExists = true
	}
	return layerExists && entityExists
}
