package data

import (
	"database/sql"
	"strconv"
)

func int64ToInt(num int64) (res int) {
	res, _ = strconv.Atoi(strconv.FormatInt(num, 10))
	return
}

func getSqlResultId(res *sql.Result) (id int) {
	id64, _ := (*res).LastInsertId()
	id = int64ToInt(id64)
	return
}
