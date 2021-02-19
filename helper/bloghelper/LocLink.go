package bloghelper

import (
	"bytes"
	"html/template"

	"rickonono3/r-blog/data"
	"rickonono3/r-blog/mytype"
)

// LocLink
//
// 描述一个归属关系中的一级目录
type LocLink struct {
	Id    int
	Title string
}

var temp = template.Must(template.New("locLink").Parse(
	`<a href="/blog/dir/{{.Id}}">{{.Title}}</a>`,
))

// MakeLocLink
//
// 将type和id确定的实体的目录归属以数组形式列出,
// 根目录为0号元素, 父目录为最后一号元素,
// 然后渲染成HTML.
func MakeLocLink(entityType, entityId int) (html template.HTML) {
	var err error
	var tmpId = entityId
	var locLinkArr = make([]LocLink, 0)
	// 如果不是目录实体, 则需要找到其父目录实体
	// 如果是目录实体, 那么链接中包含它本身, 因此不需要寻找父目录
	if entityType != 0 {
		tmpId, err = data.GetParentDir(mytype.Entity{
			Type: entityType,
			Id:   entityId,
		})
		if err != nil {
			panic("Can not get loclink")
		}
	}
	// 开始不断向上添加目录至Arr
	for {
		// 将当前级添加入Arr
		var dir mytype.Dir
		if dir, err = data.GetDir(tmpId); err != nil {
			break
		}
		locLinkArr = append([]LocLink{
			{
				Id:    tmpId,
				Title: dir.Entity.Title,
			},
		}, locLinkArr...)
		// 准备进入父级
		if tmpId == 0 {
			break
		}
		if tmpId, err = data.GetParentDir(dir.Entity); err != nil {
			break
		}
	}
	// 开始渲染HTML
	htmlStr := ""
	for i, loc := range locLinkArr {
		buf := new(bytes.Buffer)
		temp.Execute(buf, loc)
		htmlStr += buf.String()
		if i < len(locLinkArr)-1 {
			htmlStr += " &gt; "
		}
	}
	// 结束
	return template.HTML(htmlStr)
}
