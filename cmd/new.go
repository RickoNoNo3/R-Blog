package cmd

import (
	"io/ioutil"
	"os"
	"strconv"

	"../data"
	"../drawer"
	"../helper"
)

// 添加文章或目录
func New(entityType int, dirId int, inputFile *os.File) (out string) {
	if !helper.CheckEntity(data.Entity{
		Type: 1,
		Id:   dirId,
	}, true) {
		return "-1"
	}
	input, err := ioutil.ReadAll(inputFile)
	if err == nil {
		id := 0
		inputStr := string(input)
		switch entityType {
		case 0:
			id, err = data.CreateArticle(inputStr, drawer.GetTitle(inputStr), dirId)
		case 1:
			id, err = data.CreateDir(helper.CutToOneLine(inputStr), dirId)
		}
		if err == nil && id != 0 {
			out = strconv.Itoa(id)
		} else {
			out = "-1"
		}
	}
	return
}
