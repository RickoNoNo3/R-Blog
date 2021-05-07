package main

import (
	"rickonono3/r-blog/cleaner"
	"rickonono3/r-blog/data"
	"rickonono3/r-blog/helper/typehelper"
	"rickonono3/r-blog/objects"
	"rickonono3/r-blog/server"
)

// TODO: 当前版本：设置、清理工具、部署工具
// TODO: 后续版本：日志、导航、标签、搜索、打印、分享、捐赠
func main() {
	// db
	data.OpenDB("blog.db")
	defer data.CloseDB()

	// cleaner
	go cleaner.Run()
	defer cleaner.Exit()

	// server
	port := objects.Config.MustGet("ServerPort").ValInt()
	server.E.Logger.Fatal(server.E.Start(":" + typehelper.MustItoa(port)))
}
