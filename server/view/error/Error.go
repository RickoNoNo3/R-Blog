package error

import (
	"html/template"

	"github.com/labstack/echo/v4"

	"rickonono3/r-blog/helper/bloghelper"
	"rickonono3/r-blog/mytype"
)

func ViewRender(c echo.Context) (err error) {
	return c.Render(404, "error", mytype.NewGroup(mytype.ObjectList{
		"Title": mytype.NewValue(bloghelper.MakeTitle("无法找到内容")),
		"Loc":   mytype.NewValue(template.HTML(`？`)),
	}))
}
