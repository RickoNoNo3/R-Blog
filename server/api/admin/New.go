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

type newReq struct {
	Data  string `json:"data"`
	Type  int    `json:"type"`
	DirId int    `json:"dirId"`
}

type newRes struct {
	Res string `json:"res"`
	Id  int    `json:"id"`
}

// New Only for type 0 and 1
// 2 => NewResource
func New(c echo.Context) (err error) {
	req := newReq{}
	res := newRes{}
	if err = c.Bind(&req); err != nil {
		return
	}
	data.DoTx(func(tx *sqlx.Tx) (err error) {
		if datahelper.IsExists(tx, mytype.Entity{
			Type: 0,
			Id:   req.DirId,
		}) {
			req.Data = datahelper.ProcessTrim(req.Data)
			if len(req.Data) > 0 {
				switch req.Type {
				case 0:
					res.Id, err = data.CreateDir(tx, req.Data, req.DirId)
				case 1:
					title, md := datahelper.ProcessForMarkdown(req.Data)
					res.Id, err = data.CreateArticle(tx, title, md, req.DirId)
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
