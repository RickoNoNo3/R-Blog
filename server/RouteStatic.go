package server

import (
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"

	"rickonono3/r-blog/data"
	"rickonono3/r-blog/helper/datahelper"
	"rickonono3/r-blog/helper/typehelper"
	"rickonono3/r-blog/mytype"
	"rickonono3/r-blog/objects"
)

// RouteStatic
// 获取静态资源的请求在此注册:
//  - 博客文件下载: /blog/file/:id
//  - 所有文件兜底: /*
func RouteStatic() {
	// 注册所有文件的默认响应
	E.GET("/*", func(c echo.Context) error {
		dirName := objects.CWD + "public"
		path := c.Request().RequestURI
		if strings.HasSuffix(path, ".js") {
			c.Response().Header().Set("Content-Type", "text/javascript")
		} else if strings.HasSuffix(path, ".css") {
			c.Response().Header().Set("Content-Type", "text/css")
		}
		return c.File(dirName + path)
	})
	// 注册博客文件下载的响应
	E.GET("/blog/file/:id", func(c echo.Context) (err error) {
		file := mytype.File{}
		err = data.DoTx("获取文件实体", func(tx *sqlx.Tx) (err error) {
			file, err = data.GetFile(tx, typehelper.MustAtoi(c.Param("id")))
			return
		})
		if err != nil {
			return err
		}
		var (
			osInfo                       os.FileInfo
			filePath, fileName, fileSize string
			osFile                       *os.File
		)
		filePath = datahelper.GetResourcePathForServer() + datahelper.GetFileName(file.Entity.Id)
		fileName = url.PathEscape(file.Entity.Title)
		if osInfo, err = os.Stat(filePath); err == nil || os.IsExist(err) {
			fileSize = typehelper.MustItoa64(osInfo.Size())
			c.Response().Header().Set("Content-Type", "application/octet-stream")
			c.Response().Header().Set("Content-Length", fileSize)
			c.Response().Header().Set("Content-Disposition", "attachment;size="+fileSize+";filename="+fileName+";filename*="+fileName)
			if osFile, err = os.Open(filePath); err == nil {
				defer osFile.Close()
				// c.Stream和c.File因不明原因, 有一定几率出现一种在浏览器端取消下载后一直卡住, 出不来goroutine的情况, 导致文件被锁无法remove, 因此自制了一个简单的带缓冲区的响应
				buf := make([]byte, 32768)
				for {
					n, err2 := osFile.Read(buf)
					if n == 0 && err2 == io.EOF {
						break
					}
					n, err2 = c.Response().Write(buf[:n])
					if n == 0 || err2 != nil {
						err = err2
						break
					}
					c.Response().Flush()
				}
				if err == nil {
					c.Response().Status = http.StatusOK
				} else {
					// TODO(可选): 断点续传(使用content-range)
					c.Response().Status = http.StatusPartialContent
				}
			}
		}
		return
	})
}
