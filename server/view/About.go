package view

import (
	"html/template"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rickonono3/m2obj"

	"rickonono3/r-blog/helper/bloghelper"
)

func About(c echo.Context) (err error) {
	return c.Render(http.StatusOK, "about", m2obj.New(m2obj.Group{
		"Title": bloghelper.MakeTitle("关于"),
		"Loc":   template.HTML(`<a href="/about">关于</a>`),
	}))
}
