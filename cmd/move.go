package cmd

import (
	"os"
	"strconv"

	"../data"
	"../helper"
)

// 移动文章或目录
func Move(dirId int, inputFile *os.File) (out string) {
	if !helper.CheckEntity(data.Entity{
		Type: 1,
		Id:   dirId,
	}, true) {
		return "-1"
	}
	var cnt = 0
	list := helper.GetInputEntityList(inputFile)
	for _, entity := range list {
		if helper.CheckEntity(entity, false) {
			if err := data.SetEntityLayer(entity, dirId); err == nil {
				cnt++
			}
		}
	}
	out = strconv.Itoa(cnt)
	return
}
