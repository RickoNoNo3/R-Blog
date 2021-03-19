package view

import (
	"html/template"
	"net/http"

	"github.com/jmoiron/sqlx"

	"rickonono3/r-blog/data"
	"rickonono3/r-blog/helper/bloghelper"
	"rickonono3/r-blog/mytype"

	"github.com/labstack/echo/v4"
)

func Index(c echo.Context) (err error) {
	var article mytype.Article
	err = data.DoTx(func(tx *sqlx.Tx) (err error) {
		if article, err = data.GetArticle(tx, 0); err == nil {
			return c.Render(http.StatusOK, "index", mytype.NewGroup(mytype.ObjectList{
				"Title": mytype.NewValue(bloghelper.MakeTitle("")),
				"Loc":   mytype.NewValue(template.HTML(`<a href="/">主页</a>`)),
				"Intro": mytype.NewValue(article.Markdown),
			}))
		}
		return
	})
	return
}
