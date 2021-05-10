package admin

import (
	"errors"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"

	"rickonono3/r-blog/data"
	"rickonono3/r-blog/helper/datahelper"
	"rickonono3/r-blog/mytype"
)

type editReq struct {
	Data string `json:"data"`
	Type int    `json:"type"`
	Id   int    `json:"id"`
}

type editRes struct {
	Res string `json:"res"`
}

func Edit(c echo.Context) (err error) {
	req := editReq{}
	res := editRes{}
	if err = c.Bind(&req); err != nil {
		return
	}
	data.DoTx(func(tx *sqlx.Tx) (err error) {
		if datahelper.IsExists(tx, mytype.Entity{
			Type: req.Type,
			Id:   req.Id,
		}) {
			req.Data = datahelper.ProcessTrim(req.Data)
			if len(req.Data) > 0 {
				switch req.Type {
				case 0:
					err = data.EditDir(tx, req.Id, req.Data)
				case 1:
					title, md := datahelper.ProcessForMarkdown(req.Data)
					err = data.EditArticle(tx, req.Id, title, md)
				case 2:
					err = data.EditFile(tx, req.Id, req.Data)
				default:
					err = errors.New("unknown type")
				}
				if err == nil {
					res.Res = "ok"
				}
			}
		}
		return
	})
	if res.Res != "ok" {
		res.Res = "err"
	}
	return c.JSON(http.StatusOK, res)
}
