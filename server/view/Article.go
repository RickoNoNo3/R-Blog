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

func Article(c echo.Context) (err error) {
	var article mytype.Article
	articleId := typehelper.MustAtoi(c.Param("id"))
	err = data.DoTx(func(tx *sqlx.Tx) (err error) {
		if article, err = data.GetArticle(tx, articleId); err == nil {
			return c.Render(http.StatusOK, "article", mytype.NewGroup(mytype.ObjectList{
				"Title":   mytype.NewValue(bloghelper.MakeTitle(article.Entity.Title)),
				"Loc":     mytype.NewValue(bloghelper.MakeLocLink(1, articleId)),
				"Article": mytype.NewValue(article),
			}))
		}
		return
	})
	return
}
