package server

import (
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

// 获取静态资源的请求在此注册:
//  - 博客文件下载: /blog/file/:id
//  - 所有文件兜底: /*
func RouteStatic() {
	// 注册所有文件的默认响应
	E.GET("/*", func(c echo.Context) error {
		dirName := objects.Config.Get("Cwd").Val.(string) + "public"
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
		err = data.DoTx(func(tx *sqlx.Tx) (err error) {
			file, err = data.GetFile(tx, typehelper.MustAtoi(c.Param("id")))
			return
		})
		if err != nil {
			return err
		}
		var (
			osInfo                       os.FileInfo
			filePath, fileName, fileSize string
		)
		filePath = datahelper.GetFilePath(file.Entity.Id)
		fileName = url.PathEscape(file.Entity.Title)
		if osInfo, err = os.Stat(filePath); err == nil {
			fileSize = typehelper.MustItoa64(osInfo.Size())
			c.Response().Header().Set("Content-Type", "application/octet-stream")
			c.Response().Header().Set("Content-Length", fileSize)
			c.Response().Header().Set("Content-Disposition", "attachment;size="+fileSize+";filename="+fileName+";filename*="+fileName)
			return c.File(filePath)
		} else {
			return err
		}
	})
}
