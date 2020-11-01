package cmd

import (
	"../data"
	"../helper"
)

type linkResult struct {
	Link []linkResultContent `json:"link"`
}

type linkResultContent struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
}

func Link(entity data.Entity) (out string) {
	var err error
	var res linkResult
	res.Link = make([]linkResultContent, 0)
	// 闭包函数, 用来简化代码
	var addDir = func(dirId int) {
		res.Link = append([]linkResultContent{
			{
				Id:    dirId,
				Title: getDirName(dirId),
			},
		}, res.Link...)
	}
	if !helper.CheckEntity(entity, true) {
		return "-1"
	}
	// 如果是目录, 直接把自己加上
	if entity.Type == 1 {
		addDir(entity.Id)
	}
	// 开始循环向上, 直到根目录
	var nowEntity = entity
	for nowEntity.Id != 0 {
		nowId := 0
		if nowId, err = data.GetParentDir(nowEntity); err != nil {
			break
		}
		nowEntity = data.Entity{
			Type: 1,
			Id:   nowId,
		}
		addDir(nowId)
	}
	out = helper.GenerateJson(res)
	return
}

func getDirName(dirId int) (res string) {
	if dirId == 0 {
		return "博客"
	}
	dir, err := data.GetDir(dirId)
	if err != nil {
		return "未知"
	}
	return dir.Name
}
