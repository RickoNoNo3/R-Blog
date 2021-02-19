package server

import (
	"strings"

	"rickonono3/r-blog/config"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var E = echo.New()

func loadEcho() {
	E.Use(middleware.CORS())
	// TODO: git提交
	E.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Skipper: func(c echo.Context) bool {
			if strings.HasSuffix(c.Request().RequestURI, ".woff") {
				return true
			}
			return false
		},
		Level: 9,
	}))
	E.Use(echo.WrapMiddleware(M.Middleware))
	updateRender()
	staticRegister()
	viewRegister()
	apiRegister()
	if config.Get("IsInDebug").Val.(bool) {
		debugRegister()
	}
}
