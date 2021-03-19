package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
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
	switch req.Type {
	case 0:
		// TODO: dir
		res.Res = "ok"
	case 1:
		// TODO: article
		res.Res = "ok"
	case 2:
		// TODO: file
		res.Res = "ok"
	default:
		res.Res = "err"
	}
	return c.JSON(http.StatusOK, res)
}
