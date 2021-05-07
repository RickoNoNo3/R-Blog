package datahelper

import (
	"os"
	"rickonono3/r-blog/cleaner"

	"rickonono3/r-blog/helper/typehelper"
	"rickonono3/r-blog/objects"
)

// GetHashFileName 生成 hash_{hashStr}
func GetHashFileName(fileName string) string {
	hashStr := MakeHashWithStr(fileName)
	return "hash_" + hashStr[len(hashStr)-32:]
}

// GetFileName 生成 file_{fileId}
func GetFileName(fileId int) string {
	return "file_" + typehelper.MustItoa(fileId)
}

// GetResourcePathForServer 生成 {CWD}/public/resource/
func GetResourcePathForServer() string {
	path := objects.CWD + "public/resource/"
	os.MkdirAll(path, 0777)
	return path
}

// GetResourcePathForBrowser 生成 /resource/
func GetResourcePathForBrowser() string {
	return "/resource/"
}

func RemoveFileByPath(filePath string) bool {
	return cleaner.SendRequest(filePath)
}

func RemoveFileByName(fileName string) bool {
	return cleaner.SendRequest(GetResourcePathForServer() + fileName)
}
