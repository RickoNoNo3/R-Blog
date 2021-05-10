package admin

import (
	"github.com/labstack/echo/v4"
	"github.com/rickonono3/m2obj"
	"net/http"
	"rickonono3/r-blog/helper/bloghelper"
)

func Settings(c echo.Context) (err error) {
	return c.Render(http.StatusOK, "admin_settings", m2obj.New(m2obj.Group{
		"Title": bloghelper.MakeTitle("设置"),
		"Option": "settings",
	}))
}
