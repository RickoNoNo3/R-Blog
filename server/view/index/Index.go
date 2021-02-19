package index

import (
	"html/template"

	"rickonono3/r-blog/data"
	"rickonono3/r-blog/helper/bloghelper"
	"rickonono3/r-blog/mytype"

	"github.com/labstack/echo/v4"
)

func ViewRender(c echo.Context) (err error) {
	var article mytype.Article
	if article, err = data.GetArticle(0); err == nil {
		return c.Render(200, "index", mytype.NewGroup(mytype.ObjectList{
			"Title": mytype.NewValue(bloghelper.MakeTitle("")),
			"Loc":   mytype.NewValue(template.HTML(`<a href="/">主页</a>`)),
			"Intro": mytype.NewValue(article.Markdown),
		}))
	}
	return err
}
