package admin

import (
	"github.com/labstack/echo/v4"
	"github.com/rickonono3/m2obj"
	"net/http"
	"rickonono3/r-blog/logger"
	"rickonono3/r-blog/objects"
)

type settingsSaveReq struct {
	Obj m2obj.Group `json:"m2obj"`
}

type settingsSaveRes struct {
	Res string `json:"res"`
}

func SettingsSave(c echo.Context) (err error) {
	req := settingsSaveReq{}
	res := settingsSaveRes{}
	if err = c.Bind(&req); err != nil {
		return
	}
	err = objects.Config.GroupMerge(m2obj.New(req.Obj), true)
	if err == nil {
		res.Res = "ok"
	} else {
		res.Res = "err"
	}

	var op = "保存配置文件"
	if res.Res == "ok" {
		logger.L.Info("[Server]", op, res.Res)
	} else {
		logger.L.Warn("[Server]", op, err)
	}
	return c.JSON(http.StatusOK, res)
}
