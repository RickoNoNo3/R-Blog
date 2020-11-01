package cmd

import (
	"io/ioutil"
	"os"
	"strconv"

	"../data"
	"../drawer"
	"../helper"
)

// 编辑文章或目录
func Edit(entity data.Entity, inputFile *os.File) (out string) {
	if !helper.CheckEntity(entity, false) {
		return "-1"
	}
	input, err := ioutil.ReadAll(inputFile)
	if err == nil {
		inputStr := string(input)
		switch entity.Type {
		case 0:
			err = data.SetArticleMarkdownAndTitle(entity.Id, inputStr, drawer.GetTitle(inputStr))
		case 1:
			err = data.SetDirName(entity.Id, helper.CutToOneLine(inputStr))
		}
		if err == nil && entity.Id != 0 {
			out = strconv.Itoa(entity.Id)
		} else {
			out = "-1"
		}
	}
	return
}
