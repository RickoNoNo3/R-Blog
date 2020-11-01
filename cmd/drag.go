package cmd

import (
	"../data"
)

func Drag(articleId int) (out string) {
	article, err := data.GetArticle(articleId)
	if err != nil {
		return "-1"
	}
	return article.Markdown
}
