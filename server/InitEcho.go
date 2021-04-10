package server

import (
	"regexp"
	"strings"

	"rickonono3/r-blog/objects"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var E = echo.New()

func loadEcho() {
	E.Use(middleware.Logger())
	E.Use(middleware.CORS())
	E.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Skipper: func(c echo.Context) bool {
			path := c.Request().RequestURI
			// 实测.woff字体文件的压缩比大于1, 跳过Gzip
			if strings.HasSuffix(path, ".woff") {
				return true
			}
			// 供下载的文件在Gzip模式下会被切块(chunk)压缩, 造成无法指定Content-Length描述文件大小.
			// 同时还会使机器高速运转, 这在某些情况下成为传输瓶颈, 本身也对机器施压造成性能问题.
			// 当前通过如下方法识别供下载的文件:
			//   1. 请求路径以 /blog/file/ 开头
			//   2. Accept 包含 application/octet-stream
			// PS: 在响应结束后再判断的逻辑很难写, 因为Gzip是直接接管了Writer, 所以只能通过请求信息来识别.
			if strings.HasPrefix(path, "/blog/file/") ||
				regexp.MustCompile(`application/octet-stream`).MatchString(c.Request().Header.Get("Accept")) {
				return true
			}
			return false
		},
		Level: 9,
	}))
	if !objects.Config.MustGet("IsInDebug").ValBool() {
		E.Use(echo.WrapMiddleware(M.Middleware))
		updateRender()
		RouteDebug()
	}
	RouteStatic()
	RouteView()
	RouteApi()
}
