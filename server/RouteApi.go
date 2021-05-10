package server

import (
	myApi "rickonono3/r-blog/server/api"
	myApiAdmin "rickonono3/r-blog/server/api/admin"
	myMiddleware "rickonono3/r-blog/server/middleware"
)

// RouteApi
// 无视图渲染需求的(主要是返回json的)请求在此注册:
//   - /api/login
//   - /api/logout
//   - /api/myApiAdmin/new
//   - /api/myApiAdmin/newResource
//   - /api/myApiAdmin/edit
//   - /api/myApiAdmin/remove
//   - /api/myApiAdmin/move
//   - /api/myApiAdmin/settings/save
//   - /api/myApiAdmin/settings/reset
//   - /api/myApiAdmin/restart
func RouteApi() {
	// 注册api接口的响应
	api := E.Group("/api")
	api.POST("/login", myApi.Login)
	api.POST("/logout", myApi.Logout)
	// 注册api接口中管理员部分的响应
	apiAdmin := api.Group("/admin", myMiddleware.AdminAccess)
	apiAdmin.POST("/new", myApiAdmin.New)
	apiAdmin.POST("/newResource", myApiAdmin.NewResource)
	apiAdmin.POST("/edit", myApiAdmin.Edit)
	apiAdmin.POST("/remove", myApiAdmin.Remove)
	apiAdmin.POST("/move", myApiAdmin.Move)
	apiAdmin.POST("/settings/save", myApiAdmin.SettingsSave)
	apiAdmin.POST("/settings/reset", myApiAdmin.SettingsReset)
	apiAdmin.POST("/restart", myApiAdmin.Restart)
}
