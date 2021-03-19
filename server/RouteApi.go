package server

import (
	myApi "rickonono3/r-blog/server/api"
	myMiddleware "rickonono3/r-blog/server/middleware"
)

// 无视图渲染需求的(主要是返回json的)请求在此注册:
//   - /api/login
//   - /api/logout
// TODO: /api/admin/new
// TODO: /api/admin/newResource
// TODO: /api/admin/edit
// TODO: /api/admin/remove
// TODO: /api/admin/move
//
func RouteApi() {
	// 注册api接口的响应
	api := E.Group("/api")
	api.POST("/login", myApi.Login)
	api.POST("/logout", myApi.Logout)
	// 注册api接口中管理员部分的响应
	apiAdmin := api.Group("/admin", myMiddleware.AdminAccess)
	apiAdmin.POST("/new", myApi.New)
	apiAdmin.POST("/newResource", myApi.NewResource)
	apiAdmin.POST("/edit", myApi.Edit)
}
