package api

import (
	"net/http"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"

	"rickonono3/r-blog/data"
	"rickonono3/r-blog/helper/datahelper"
	"rickonono3/r-blog/mytype"
)

type removeReqItem struct {
	Type int `json:"type"`
	Id   int `json:"id"`
}

type removeReq struct {
	List []removeReqItem `json:"list"`
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
		if err = data.DoTx(func(tx *sqlx.Tx) (err error) {
			entity := mytype.Entity{
				Type: item.Type,
				Id:   item.Id,
			}
			if datahelper.IsExists(tx, entity) {
				if item.Type == 2 {
					filePath := datahelper.GetResourcePathForServer() + datahelper.GetFileName(item.Id)
					_, err = os.Stat(filePath)
					if os.IsNotExist(err) {
						err = nil
					} else if err == nil {
						err = os.Remove(filePath)
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
	return c.JSON(http.StatusOK, res)
}
