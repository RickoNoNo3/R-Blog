package cmd

import (
	"strconv"

	"../data"
	"../drawer"
	"../helper"
)

// 带数据库操作的渲染, 只给定文章Id
func Draw(id int) (out string) {
	if !helper.CheckEntity(data.Entity{
		Type: 0,
		Id:   id,
	}, false) {
		return "-1"
	}
	var htmlres string
	article, err := data.GetArticle(id) // 从数据库获取id对应的md
	if err != nil {
		return "-1"
	}
	_, htmlres = drawer.Draw(article.Markdown) // 渲染
	err = data.SetArticleHtml(id, htmlres)     // 把渲染结果放回数据库
	if err != nil {
		return "-1"
	}
	out = strconv.Itoa(id)
	return
}
