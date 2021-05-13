package admin

import (
	"encoding/binary"
	"io"
	"net/http"
	"net/url"
	"os"
	"rickonono3/r-blog/logger"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"

	"rickonono3/r-blog/data"
	"rickonono3/r-blog/helper/datahelper"
	"rickonono3/r-blog/helper/typehelper"
)

//   4           FileNameLen 1     [4]     n
//   FileNameLen FileName    IsTmp [DirId] FileData
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
	var filePath string
	err = data.DoTx(func(tx *sqlx.Tx) (err error) {
		// 读取二进制body
		if bodyReader := c.Request().Body; bodyReader != nil && c.Request().ContentLength >= 8 {
			// 读取4字节 FileNameLen
			req.FileNameLen = binary.BigEndian.Uint32(readNBytes(bodyReader, 4))
			// 读取 FileNameLen 字节 FileName
			req.FileName, _ = url.PathUnescape(strings.TrimSpace(string(readNBytes(bodyReader, req.FileNameLen))))
			// 若FileName有效则继续
			if req.FileName != "" {
				// 读取1字节 IsTmp
				req.IsTmp = readNBytes(bodyReader, 1)[0]
				if req.IsTmp == 1 { // 是临时文件, 开始写文件吧
					fileName := datahelper.GetHashFileName(req.FileName)
					filePath = datahelper.GetResourcePathForServer() + fileName
					if err = writeToFile(
						bodyReader,
						filePath,
					); err == nil {
						res.Res = "ok"
						res.FileLoc = datahelper.GetResourcePathForBrowser() + fileName
					} else {
						datahelper.RemoveFileByPath(filePath)
					}
				} else { // 是固定文件, 先创建数据库再写入文件
					// 读取4字节 DirId
					req.DirId = int(binary.BigEndian.Uint32(readNBytes(bodyReader, 4)))
					// 在数据库中创建文件索引
					fileId := 0
					if fileId, err = data.CreateFile(tx, req.FileName, req.DirId); err == nil {
						// 开始写文件
						fileName := datahelper.GetFileName(fileId)
						filePath = datahelper.GetResourcePathForServer() + fileName
						if err = writeToFile(
							bodyReader,
							filePath,
						); err == nil {
							res.Res = "ok"
							res.FileLoc = "/blog/file/" + typehelper.MustItoa(fileId)
						} else {
							datahelper.RemoveFileByPath(filePath)
						}
					}
				}
			}
		}
		return
	})
	if res.Res != "ok" {
		res.Res = "err"
	}

	var op = strings.Join([]string{
		"新建",
		datahelper.GetTypeName(2),
		"到",
		filePath,
		": ",
	}, "")
	if res.Res == "ok" {
		logger.L.Info("[Server]", op, res.Res)
	} else {
		logger.L.Warn("[Server]", op, err)
	}

	if err != nil {
		return err
	} else {
		return c.JSON(http.StatusOK, res)
	}
}

// 从当前指针开始读取 n 字节从 reader 中
func readNBytes(reader io.Reader, n uint32) []byte {
	buf := make([]byte, n)
	io.ReadFull(reader, buf)
	return buf
}

// 把 reader 内剩余的所有字节写入 filePath 文件
func writeToFile(reader io.Reader, filePath string) (err error) {
	var file *os.File
	if file, err = os.Create(filePath); err == nil {
		defer file.Close()
		if _, err = io.Copy(file, reader); err == nil {
			err = file.Close()
		}
	}
	return
}
