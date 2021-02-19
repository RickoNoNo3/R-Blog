package server

import (
	"strings"

	"github.com/labstack/echo/v4"

	"rickonono3/r-blog/config"
)

func staticRegister() {
	dirName := config.Get("Cwd").Val.(string) + "public"
	E.GET("/*", func(c echo.Context) error {
		path := c.Request().RequestURI
		if strings.HasSuffix(path, ".js") {
			c.Response().Header().Set("Content-Type", "text/javascript")
		} else if strings.HasSuffix(path, ".css") {
			c.Response().Header().Set("Content-Type", "text/css")
		}
		return c.File(dirName + path)
	})
}
