package datahelper

import (
	"rickonono3/r-blog/helper/typehelper"
	"rickonono3/r-blog/objects"
)

func GetFilePath(fileId int) string {
	return objects.Config.Get("Cwd").Val.(string) + GetFilePathAbsolutely(fileId)
}

func GetResourcePath() string {
	return objects.Config.Get("Cwd").Val.(string) + GetResourcePathAbsolutely()
}

func GetFilePathAbsolutely(fileId int) string {
	return GetResourcePathAbsolutely() + "file_" + typehelper.MustItoa(fileId)
}

func GetResourcePathAbsolutely() string {
	return "public/resource/"
}
