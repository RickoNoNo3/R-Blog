package server

import (
	"regexp"

	"rickonono3/r-blog/config"
	"rickonono3/r-blog/server/view/article"
	"rickonono3/r-blog/server/view/dir"
	myError "rickonono3/r-blog/server/view/error"
	"rickonono3/r-blog/server/view/index"

	"github.com/labstack/echo/v4"
)

// 在全部view中都套用此高阶函数, 便于模板热更(限DEBUG模式)
var viewContainer = func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if config.Get("IsInDebug").Val.(bool) {
			updateRender()
		}
		return next(c)
	}
}

func viewRegister() {
	E.HTTPErrorHandler = func(err error, c echo.Context) {
		accept := c.Request().Header.Get("Accept")
		path := c.Request().RequestURI
		reg1 := regexp.MustCompile(`application/(json|octet-stream)`)
		reg2 := regexp.MustCompile(`^/api/`)
		if c.Request().Method == "GET" && !reg1.MatchString(accept) && !reg2.MatchString(path) {
			err = viewContainer(myError.ViewRender)(c)
		} else {
			E.DefaultHTTPErrorHandler(err, c)
		}
	}
	E.GET("/", viewContainer(index.ViewRender))
	E.GET("/blog/dir/:id", viewContainer(dir.ViewRender))
	E.GET("/blog/article/:id", viewContainer(article.ViewRender))
	// TODO: File
	// E.GET("/blog/file/:id", viewContainer(file.ViewRender))
	E.GET("/blog/*", func(c echo.Context) error {
		return c.Redirect(301, "/blog/dir/0")
	})
}
