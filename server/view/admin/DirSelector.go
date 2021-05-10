package admin

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rickonono3/m2obj"

	"rickonono3/r-blog/helper/datahelper"
)

func DirSelector(c echo.Context) (err error) {
	var list []datahelper.DirInfo
	list, err = datahelper.GetDirInfoList()
	if err != nil {
		return
	}
	return c.Render(http.StatusOK, "dir_selector", m2obj.New(m2obj.Group{
		"List": list,
	}))
}
