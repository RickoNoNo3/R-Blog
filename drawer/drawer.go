package drawer

// Drawer 通用的Draw功能: Parse和Render
// 输入markdown字符串, 输出渲染好的html页面和标题
func Draw(md string) (title, htmlres string) {
	title, htmlres = Parse(md)
	htmlres = Render(htmlres)
	return
}

// 获取md的title
func GetTitle(md string) (title string) {
	md = makeMd(md)
	title = makeTitle(md)
	return
}
