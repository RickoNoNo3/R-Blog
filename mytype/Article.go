package mytype

import "strings"

type Article struct {
	Entity   Entity
	Markdown string
	Tags     string
	Voted    int
	Visited  int
}

func (article *Article) GetTags() (tagsList []string) {
	tagsList = make([]string, 0)
	if tags := strings.TrimSpace(article.Tags); tags != "" {
		tagsList = append(tagsList, strings.Split(tags, "|")...)
	}
	return
}
