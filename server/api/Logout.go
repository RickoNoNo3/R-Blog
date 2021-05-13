package api

import (
	"net/http"
	"rickonono3/r-blog/logger"
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
	logger.L.Info("[Server]", "管理员注销: ", c.RealIP())
	hash := datahelper.MakeHashWithStr("")
	objects.RuntimeEnv.Set("AdminHash", hash)
	c.SetCookie(&http.Cookie{
		Name:    "blog-login",
		Value:   "",
		Expires: time.Now().Add(-1 * time.Hour),
		Path:    "/",
	})
	res.Res = "ok"
	return c.JSON(http.StatusOK, res)
}
