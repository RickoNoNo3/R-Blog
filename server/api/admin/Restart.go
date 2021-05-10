package admin

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"rickonono3/r-blog/helper/cmdhelper"
)

type restartRes struct {
	Res string `json:"res"`
}

func Restart(c echo.Context) (err error) {
	cmdhelper.Input("restart")
	return c.JSON(http.StatusOK, restartRes{
		Res: "ok",
	})
}
