package api

import (
	"encoding/binary"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"

	"rickonono3/r-blog/data"
	"rickonono3/r-blog/helper/datahelper"
)

type newResourceReq struct {
	FileNameLen uint32
	FileName    string
	IsTmp       byte
	DirId       int
}

type newResourceRes struct {
	Res     string `json:"res"`
	FileLoc string `json:"fileLoc"`
}

func NewResource(c echo.Context) (err error) {
	req := newResourceReq{}
	res := newResourceRes{}
	err = data.DoTx(func(tx *sqlx.Tx) (err error) {
		if bodyReader := c.Request().Body; bodyReader != nil && c.Request().ContentLength >= 8 {
			req.FileNameLen = binary.BigEndian.Uint32(readNBytes(bodyReader, 4))
			req.FileName, _ = url.PathUnescape(strings.TrimSpace(string(readNBytes(bodyReader, req.FileNameLen))))
			if req.FileName != "" {
				req.IsTmp = readNBytes(bodyReader, 1)[0]
				if req.IsTmp == 1 { // 是临时文件
					fileHash := datahelper.MakeHashWithStr(req.FileName)
					fileLoc := datahelper.GetResourcePath() + fileHash
					if err = writeToFile(bodyReader, fileLoc); err == nil {
						res.Res = "ok"
						res.FileLoc = datahelper.GetResourcePathAbsolutely() + fileHash
						return c.JSON(http.StatusOK, res)
					}
				} else { // 是固定文件
					req.DirId = int(binary.BigEndian.Uint32(readNBytes(bodyReader, 4)))
					fileId := 0
					if fileId, err = data.CreateFile(tx, req.FileName, req.DirId); err == nil {
						fileLoc := datahelper.GetFilePath(fileId)
						if err = writeToFile(bodyReader, fileLoc); err == nil {
							res.Res = "ok"
							return c.JSON(http.StatusOK, res)
						}
					}
				}
			}
			// 处理好前面的数据, 剩下的就是FileData了, 直接io.Copy
		}
		return
	})
	if err != nil {
		res.Res = "err"
		return c.JSON(http.StatusOK, res)
	} else {
		return
	}
}

func readNBytes(reader io.Reader, n uint32) []byte {
	buf := make([]byte, n)
	io.ReadFull(reader, buf)
	return buf
}

func writeToFile(reader io.Reader, fileLoc string) (err error) {
	var file *os.File
	if file, err = os.Create(fileLoc); err == nil {
		if _, err = io.Copy(file, reader); err == nil {
			file.Close()
		}
	}
	return
}
