package server

import (
	myApi "rickonono3/r-blog/server/api"
	myMiddleware "rickonono3/r-blog/server/middleware"
)

// RouteApi
// 无视图渲染需求的(主要是返回json的)请求在此注册:
//   - /api/login
//   - /api/logout
//   - /api/admin/new
//   - /api/admin/newResource
//   - /api/admin/edit
//   - /api/admin/remove
//   - /api/admin/move
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
	apiAdmin.POST("/remove", myApi.Remove)
	apiAdmin.POST("/move", myApi.Move)
}
