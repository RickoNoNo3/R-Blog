package api

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"rickonono3/r-blog/helper/bloghelper"
	"rickonono3/r-blog/mytype"
	"rickonono3/r-blog/objects"
)

type logoutRes struct {
	Res string `json:"res"`
}

func Logout(c echo.Context) (err error) {
	res := logoutRes{}
	hash := bloghelper.GetAdminHash()
	objects.Cache.Set("AdminHash", mytype.NewValue(hash))
	c.SetCookie(&http.Cookie{
		Name:    "blog-login",
		Value:   "",
		Expires: time.Now().Add(-1 * time.Hour),
		Path:    "/",
	})
	res.Res = "ok"
	return c.JSON(http.StatusOK, res)
}
