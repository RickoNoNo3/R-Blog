package main

import (
	"rickonono3/r-blog/server"
)

func main() {
	server.E.Logger.Fatal(server.E.Start(":13808"))
}
