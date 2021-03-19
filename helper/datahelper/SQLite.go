package datahelper

import (
	"database/sql"
	"rickonono3/r-blog/helper/typehelper"
)

func GetResultId(res *sql.Result) (id int) {
	id64, _ := (*res).LastInsertId()
	id = typehelper.Int64ToInt(id64)
	return
}
