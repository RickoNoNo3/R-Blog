package datahelper

import (
	"rickonono3/r-blog/helper/typehelper"
	"rickonono3/r-blog/mytype"
)

func GetTypeName(entityType int) string {
	switch entityType {
	case 0:
		return "目录"
	case 1:
		return "文章"
	case 2:
		return "文件"
	}
	return "未知内容"
}

func GetEntityStr(entity mytype.EasyEntity) string {
	str := GetTypeName(entity.Type) + "-"
	str += typehelper.MustItoa(entity.Id)
	return str
}

func GetEntityListStr(entityList []mytype.EasyEntity) []string {
	res := make([]string, len(entityList))
	for i, entity := range entityList {
		res[i] = GetEntityStr(entity)
	}
	return res
}
