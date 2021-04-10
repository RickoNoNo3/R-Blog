package view

import (
	"html/template"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rickonono3/m2obj"

	"rickonono3/r-blog/helper/bloghelper"
)

func Error(c echo.Context) (err error) {
	return c.Render(http.StatusNotFound, "error", m2obj.New(m2obj.Group{
		"Title": bloghelper.MakeTitle("无法找到内容"),
		"Loc":   template.HTML(`？`),
	}))
}
