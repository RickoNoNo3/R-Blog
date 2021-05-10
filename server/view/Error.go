package view

import (
	"github.com/labstack/echo/v4"
	"github.com/rickonono3/m2obj"
	"html/template"

	"rickonono3/r-blog/helper/bloghelper"
)

func ErrorWrap(statusCode int) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		return c.Render(statusCode, "error", m2obj.New(m2obj.Group{
			"Title": bloghelper.MakeTitle("无法找到内容"),
			"Loc":   template.HTML(`？`),
		}))
	}
}
