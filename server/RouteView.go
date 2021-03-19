package server

import (
	"net/http"
	"regexp"

	"rickonono3/r-blog/objects"
	myMiddleware "rickonono3/r-blog/server/middleware"
	"rickonono3/r-blog/server/view"
	viewAdmin "rickonono3/r-blog/server/view/admin"

	"github.com/labstack/echo/v4"
)

// viewContainer 在全部view中都套用此高阶函数, 便于模板热更(限DEBUG模式)
func viewContainer(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if objects.Config.Get("IsInDebug").Val.(bool) {
			updateRender()
		}
		return next(c)
	}
}

// 需要渲染视图的请求在此注册:
//  - err时的响应
//  - 主页: /
//  - 博客路径重定向: /blog/
//  - 博客目录页: /blog/dir/:id
//  - 博客文章页: /blog/article/:id
//  - 管理员路径重定向: /admin/
//  - 管理员编辑: /admin/edit?type&id&parentId
// TODO: /admin/tool/dirSelector
// TODO: /admin/settings
func RouteView() {
	// 注册错误处理器, 处理发生错误时的响应
	E.HTTPErrorHandler = func(err error, c echo.Context) {
		accept := c.Request().Header.Get("Accept")
		path := c.Request().RequestURI
		reg1 := regexp.MustCompile(`application/(json|octet-stream)`)
		reg2 := regexp.MustCompile(`^/api/`)
		if c.Request().Method == "GET" && !reg1.MatchString(accept) && !reg2.MatchString(path) {
			err = viewContainer(view.Error)(c)
		} else {
			E.DefaultHTTPErrorHandler(err, c)
		}
	}
	// 注册主页与博客页面
	E.GET("/", viewContainer(view.Index))
	E.GET("/blog/", func(c echo.Context) error {
		return c.Redirect(http.StatusPermanentRedirect, "/blog/dir/0")
	})
	E.GET("/blog/dir/:id", viewContainer(view.Dir))
	E.GET("/blog/article/:id", viewContainer(view.Article))
	// 注册管理员页面
	admin := E.Group("/admin", myMiddleware.AdminAccess)
	admin.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusPermanentRedirect, "/admin/edit")
	})
	admin.GET("/edit", viewContainer(viewAdmin.Edit))
}
