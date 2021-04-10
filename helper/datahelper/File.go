package datahelper

import (
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
	return objects.Config.MustGet("Cwd").ValStr() + "public/resource/"
}

// 生成 /resource/
func GetResourcePathForBrowser() string {
	return "/resource/"
}
