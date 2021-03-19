package main

import (
	"rickonono3/r-blog/data"
	"rickonono3/r-blog/server"
)

func main() {
	defer data.CloseDB()
	data.OpenDB("blog.db")
	server.E.Logger.Fatal(server.E.Start(":13808"))
}
