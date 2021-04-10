package api

import (
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"

	"rickonono3/r-blog/data"
	"rickonono3/r-blog/helper/datahelper"
	"rickonono3/r-blog/mytype"
)

type moveReqItem struct {
	Type int `json:"type"`
	Id   int `json:"id"`
}

type moveReq struct {
	List  []moveReqItem `json:"list"`
	DirId int           `json:"dirId"`
}

type moveRes struct {
	Res string `json:"res"`
}

func Move(c echo.Context) (err error) {
	req := moveReq{}
	res := moveRes{}
	if err = c.Bind(&req); err != nil {
		return
	}
	dirExists := false
	data.DoTx(func(tx *sqlx.Tx) (err error) {
		dirExists = datahelper.IsExists(tx, mytype.Entity{
			Type: 0,
			Id:   req.DirId,
		})
		return
	})
	for _, item := range req.List {
		if err = data.DoTx(func(tx *sqlx.Tx) (err error) {
			entity := mytype.Entity{
				Type: item.Type,
				Id:   item.Id,
			}
			if datahelper.IsExists(tx, entity) {
				err = data.MoveLayer(tx, entity, req.DirId)
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