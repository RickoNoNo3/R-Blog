package admin

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"rickonono3/r-blog/logger"
	"rickonono3/r-blog/objects"
)

type settingsResetRes struct {
	Res string `json:"res"`
}

func SettingsReset(c echo.Context) (err error) {
	logger.L.Warn("[Server]", "重置配置文件")
	objects.Config.SetVal(objects.DefaultConfig.Clone())
	return c.JSON(http.StatusOK, settingsResetRes{
		Res: "ok",
	})
}
