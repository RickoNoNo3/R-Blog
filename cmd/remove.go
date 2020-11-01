package cmd

import (
	"os"
	"strconv"

	"../data"
	"../helper"
)

// 删除文章或目录
func Remove(inputFile *os.File) (out string) {
	var err error
	var cnt = 0
	list := helper.GetInputEntityList(inputFile)
	for _, entity := range list {
		if helper.CheckEntity(entity, false) {
			switch entity.Type {
			case 0:
				err = data.DeleteArticle(entity.Id)
			case 1:
				err = data.DeleteDir(entity.Id)
			}
			if err == nil {
				cnt++
			}
		}
	}
	out = strconv.Itoa(cnt)
	return
}
