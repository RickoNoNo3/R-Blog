package article

import (
	"github.com/labstack/echo/v4"

	"rickonono3/r-blog/data"
	"rickonono3/r-blog/helper/bloghelper"
	"rickonono3/r-blog/helper/typehelper"
	"rickonono3/r-blog/mytype"
)

func ViewRender(c echo.Context) (err error) {
	var article mytype.Article
	articleId := typehelper.MustAtoi(c.Param("id"))
	if article, err = data.GetArticle(articleId); err == nil {
		return c.Render(200, "article", mytype.NewGroup(mytype.ObjectList{
			"Title":   mytype.NewValue(bloghelper.MakeTitle(article.Entity.Title)),
			"Loc":     mytype.NewValue(bloghelper.MakeLocLink(0, articleId)),
			"Article": mytype.NewValue(article),
		}))
	}
	return err
}
