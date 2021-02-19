package dir

import (
	"github.com/labstack/echo/v4"

	"rickonono3/r-blog/data"
	"rickonono3/r-blog/helper/bloghelper"
	"rickonono3/r-blog/helper/typehelper"
	"rickonono3/r-blog/mytype"
)

func ViewRender(c echo.Context) (err error) {
	var dir mytype.Dir
	var contents []mytype.Entity
	dirId := typehelper.MustAtoi(c.Param("id"))
	if dir, err = data.GetDir(dirId); err == nil {
		if contents, err = data.GetContents(dirId); err == nil {
			return c.Render(200, "dir", mytype.NewGroup(mytype.ObjectList{
				"Title":    mytype.NewValue(bloghelper.MakeTitle(dir.Entity.Title)),
				"Loc":      mytype.NewValue(bloghelper.MakeLocLink(0, dirId)),
				"Dir":      mytype.NewValue(dir),
				"Contents": mytype.NewValue(contents),
			}))
		}
	}
	return err
}
