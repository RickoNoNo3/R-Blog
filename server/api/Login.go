package api

import (
	"net/http"
	"rickonono3/r-blog/logger"
	"time"

	"github.com/labstack/echo/v4"

	"rickonono3/r-blog/helper/datahelper"
	"rickonono3/r-blog/objects"
)

type loginReq struct {
	Pswd string `json:"pswd"`
}

type loginRes struct {
	Res string `json:"res"`
}

func Login(c echo.Context) (err error) {
	req := loginReq{}
	res := loginRes{}
	if err = c.Bind(&req); err != nil {
		return
	}
	if req.Pswd == objects.Config.MustGet("AdminPSWD").ValStr() {
		logger.L.Info("[Server]", "管理员登录: ", c.RealIP())
		hash := datahelper.MakeHashWithStr(c.RealIP())
		objects.RuntimeEnv.Set("AdminHash", hash)
		c.SetCookie(&http.Cookie{
			Name:    "blog-login",
			Value:   hash,
			Expires: time.Now().Add(2 * time.Hour),
			Path:    "/",
		})
		res.Res = "ok"
	} else {
		res.Res = "err"
	}
	return c.JSON(http.StatusOK, res)
}
