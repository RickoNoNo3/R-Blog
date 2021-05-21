package view

import (
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/rickonono3/m2obj"

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
	err = data.DoTx("获取目录实体", func(tx *sqlx.Tx) (err error) {
		if dir, err = data.GetDir(tx, dirId); err == nil {
			var contents []mytype.Entity
			if contents, err = data.GetContents(tx, dirId); err == nil {
				return c.Render(http.StatusOK, "dir", m2obj.New(m2obj.Group{
					"Title":    bloghelper.MakeTitle(dir.Entity.Title),
					"Loc":      bloghelper.MakeLocLink(0, dirId),
					"Dir":      dir,
					"Contents": contents,
				}))
			}
		}
		return
	})
	return
}
