package view

import (
	"html/template"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/rickonono3/m2obj"

	"github.com/labstack/echo/v4"
	"rickonono3/r-blog/data"
	"rickonono3/r-blog/helper/bloghelper"
)

func Index(c echo.Context) (err error) {
	var welcome string
	err = data.DoTx(func(tx *sqlx.Tx) (err error) {
		if welcome, err = data.GetWelcome(tx); err == nil {
			return c.Render(http.StatusOK, "index", m2obj.New(m2obj.Group{
				"Title": bloghelper.MakeTitle(""),
				"Loc":   template.HTML(`<a href="/">主页</a>`),
				"Intro": welcome,
			}))
		}
		return
	})
	return
}
