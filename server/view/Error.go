package view

import (
	"html/template"
	"net/http"

	"github.com/labstack/echo/v4"

	"rickonono3/r-blog/helper/bloghelper"
	"rickonono3/r-blog/mytype"
)

func Error(c echo.Context) (err error) {
	return c.Render(http.StatusNotFound, "error", mytype.NewGroup(mytype.ObjectList{
		"Title": mytype.NewValue(bloghelper.MakeTitle("无法找到内容")),
		"Loc":   mytype.NewValue(template.HTML(`？`)),
	}))
}
