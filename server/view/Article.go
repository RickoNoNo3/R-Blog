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

func Article(c echo.Context) (err error) {
	var article mytype.Article
	articleId := typehelper.MustAtoi(c.Param("id"))
	err = data.DoTx("获取文章实体", func(tx *sqlx.Tx) (err error) {
		if article, err = data.GetArticle(tx, articleId); err == nil {
			return c.Render(http.StatusOK, "article", m2obj.New(m2obj.Group{
				"Title":   bloghelper.MakeTitle(article.Entity.Title),
				"Loc":     bloghelper.MakeLocLink(1, articleId),
				"Article": article,
			}))
		}
		return
	})
	return
}
