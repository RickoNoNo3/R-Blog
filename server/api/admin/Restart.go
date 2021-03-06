package admin

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"rickonono3/r-blog/helper/cmdhelper"
	"rickonono3/r-blog/logger"
)

type restartRes struct {
	Res string `json:"res"`
}

func Restart(c echo.Context) (err error) {
	logger.L.Warn("[Server]", "重启服务器")
	cmdhelper.Input("restart")
	return c.JSON(http.StatusOK, restartRes{
		Res: "ok",
	})
}
