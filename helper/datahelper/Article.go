package datahelper

import (
	"regexp"
	"strings"
)

func ProcessTrim(str string) string {
	return strings.Trim(str, "\n\r 　\t")
}

func ProcessForMarkdown(md string) (title, md2 string) {
	// 处理md
	{
		md = strings.ReplaceAll(md, "\r\n", "\n")
		// 将md内的html标签全转换成小写
		regexp.MustCompile(`<[\S\s]+?>`).ReplaceAllStringFunc(md, strings.ToLower)
	}

	// 从第一行中截取出title
	{
		firstLine := strings.TrimSpace(strings.Split(md, "\n")[0])
		titleBeg, titleEnd := strings.Index(firstLine, "#"), strings.Index(md, "\n")
		// 就一行, 特判
		if titleEnd == -1 {
			titleEnd = len(firstLine)
		}
		// 如果第一行没井号就当作第一行全部为标题
		if titleBeg != -1 {
			title = strings.TrimSpace(firstLine[titleBeg+1 : titleEnd])
		} else {
			title = firstLine
		}
	}
	md2 = md
	return
}
