package datahelper

import (
	"os"

	"rickonono3/r-blog/helper/typehelper"
	"rickonono3/r-blog/objects"
)

// 生成 hash_{hashStr}
func GetHashFileName(fileName string) string {
	hashStr := MakeHashWithStr(fileName)
	return "hash_" + hashStr[len(hashStr)-32:]
}

// 生成 file_{fileId}
func GetFileName(fileId int) string {
	return "file_" + typehelper.MustItoa(fileId)
}

// 生成 {Cwd}/public/resource/
func GetResourcePathForServer() string {
	path := objects.Config.MustGet("Cwd").ValStr() + "public/resource/"
	os.MkdirAll(path, 0777)
	return path
}

// 生成 /resource/
func GetResourcePathForBrowser() string {
	return "/resource/"
}
