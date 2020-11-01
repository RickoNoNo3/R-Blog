package cmd

import (
	"../data"
	"../helper"
)

type articleResult struct {
	Title     string `json:"title"`
	Html      string `json:"html"`
	CreatedT  string `json:"createdT"`
	ModifiedT string `json:"modifiedT"`
	// tags       []string
	// contents   []{
	//   text       string
	//   id         string
	//   contents   []{...}
	// }
}

type dirResult struct {
	List []dirResultContent `json:"list"`
}

type dirResultContent struct {
	Id        int    `json:"id"`
	Type      int    `json:"type"`
	Text      string `json:"text"`
	CreatedT  string `json:"createdT"`
	ModifiedT string `json:"modifiedT"`
}

func Read(entity data.Entity) (out string) {
	if !helper.CheckEntity(entity, true) {
		return "-1"
	}
	var res interface{}
	switch entity.Type {
	case 0:
		res = readArticle(&entity)
	case 1:
		res = readDir(&entity)
	}
	if res != nil {
		out = helper.GenerateJson(res)
	}
	return
}

func readArticle(entity *data.Entity) (res articleResult) {
	drewT, createdT, modifiedT, err := data.GetArticleTime(entity.Id)
	if err == nil {
		// 添加字段
		res.CreatedT = createdT
		res.ModifiedT = modifiedT
		// 修改比渲染新, 需要重新渲染
		drew, modified := helper.ParseStringToTime(drewT), helper.ParseStringToTime(modifiedT)
		if modified.After(drew) {
			Draw(entity.Id)
		}
		// 添加字段
		article, err := data.GetArticle(entity.Id)
		if err == nil {
			res.Title = article.Title
			res.Html = article.Html
		}
	}
	return
}

func readDir(entity *data.Entity) (res dirResult) {
	res.List = make([]dirResultContent, 0)
	contents, err := data.GetLayer(entity.Id)
	if err == nil {
		for _, content := range contents {
			res.List = append(res.List, dirResultContent{
				Id:        content.Id,
				Type:      content.Type,
				Text:      content.Text,
				CreatedT:  content.CreatedT,
				ModifiedT: content.ModifiedT,
			})
		}
	}
	return
}
