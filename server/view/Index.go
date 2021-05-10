package view

import (
	"html/template"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/rickonono3/m2obj"

	"rickonono3/r-blog/data"
	"rickonono3/r-blog/helper/bloghelper"
	"rickonono3/r-blog/mytype"

	"github.com/labstack/echo/v4"
)

func Index(c echo.Context) (err error) {
	var article mytype.Article
	err = data.DoTx(func(tx *sqlx.Tx) (err error) {
		if article, err = data.GetArticle(tx, 0); err == nil {
			return c.Render(http.StatusOK, "index", m2obj.New(m2obj.Group{
				"Title": bloghelper.MakeTitle(""),
				"Loc":   template.HTML(`<a href="/">主页</a>`),
				"Intro": article.Markdown,
			}))
		}
		return
	})
	return
}
