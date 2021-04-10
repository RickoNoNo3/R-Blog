package api

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"rickonono3/r-blog/helper/datahelper"
	"rickonono3/r-blog/objects"
)

type logoutRes struct {
	Res string `json:"res"`
}

func Logout(c echo.Context) (err error) {
	res := logoutRes{}
	hash := datahelper.MakeHashWithStr("")
	objects.Cache.Set("AdminHash", hash)
	c.SetCookie(&http.Cookie{
		Name:    "blog-login",
		Value:   "",
		Expires: time.Now().Add(-1 * time.Hour),
		Path:    "/",
	})
	res.Res = "ok"
	return c.JSON(http.StatusOK, res)
}
