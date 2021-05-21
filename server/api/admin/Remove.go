package admin

import (
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"net/http"
	"rickonono3/r-blog/logger"
	"strings"

	"rickonono3/r-blog/data"
	"rickonono3/r-blog/helper/datahelper"
	"rickonono3/r-blog/mytype"
)

type removeReq struct {
	List []mytype.EasyEntity `json:"list"`
}

type removeRes struct {
	Res string `json:"res"`
}

func Remove(c echo.Context) (err error) {
	req := removeReq{}
	res := removeRes{}
	if err = c.Bind(&req); err != nil {
		return
	}
	for _, item := range req.List {
		if err = data.DoTx("删除实体", func(tx *sqlx.Tx) (err error) {
			entity := mytype.Entity{
				Type: item.Type,
				Id:   item.Id,
			}
			if datahelper.IsExists(tx, entity) {
				if item.Type == 2 {
					if !datahelper.RemoveFileByName(datahelper.GetFileName(item.Id)) {
						err = errors.New("removing file failed - cleaning goroutine is exited")
					}
				}
				if err == nil {
					err = data.RemoveLayer(tx, entity)
				}
			}
			return
		}); err != nil {
			break
		}
	}
	if err == nil {
		res.Res = "ok"
	} else {
		res.Res = "err"
	}

	var op = strings.Join([]string{
		"删除[",
		strings.Join(datahelper.GetEntityListStr(req.List), ","),
		"]: ",
	}, "")
	if res.Res == "ok" {
		logger.L.Info("[Server]", op, res.Res)
	} else {
		logger.L.Warn("[Server]", op, err)
	}
	return c.JSON(http.StatusOK, res)
}
