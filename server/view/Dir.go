package view

import (
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"

	"rickonono3/r-blog/data"
	"rickonono3/r-blog/helper/bloghelper"
	"rickonono3/r-blog/helper/typehelper"
	"rickonono3/r-blog/mytype"
)

// 专属参数:
//   Dir
//   Contents

func Dir(c echo.Context) (err error) {
	dirId := typehelper.MustAtoi(c.Param("id"))
	var dir mytype.Dir
	err = data.DoTx(func(tx *sqlx.Tx) (err error) {
		if dir, err = data.GetDir(tx, dirId); err == nil {
			var contents []mytype.Entity
			if contents, err = data.GetContents(tx, dirId); err == nil {
				return c.Render(http.StatusOK, "dir", mytype.NewGroup(mytype.ObjectList{
					"Title":    mytype.NewValue(bloghelper.MakeTitle(dir.Entity.Title)),
					"Loc":      mytype.NewValue(bloghelper.MakeLocLink(0, dirId)),
					"Dir":      mytype.NewValue(dir),
					"Contents": mytype.NewValue(contents),
				}))
			}
		}
		return
	})
	return
}
