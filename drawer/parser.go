package drawer

import (
	"regexp"
	"strings"

	"gopkg.in/russross/blackfriday.v2"
)

// Parse 把markdown解析为html. 是很原始的markdown解析
// JS插件引擎需要进一步通过Render进行渲染
func Parse(md string) (title, htmlres string) {
	md = makeMd(md)
	title = makeTitle(md)
	htmlres = makeHtml(md)
	return
}

// 对md进行最简单的格式处理
func makeMd(md string) string {
	md = strings.ReplaceAll(md, "\r\n", "\n")
	// 将md内的html标签全转换成小写
	reg, _ := regexp.Compile(`<[\S\s]+?>`)
	md = reg.ReplaceAllStringFunc(md, strings.ToLower)
	// trim
	md = strings.Trim(md, "\n 　\t")
	// fmt.Println(md)
	return md
}

// 拿标题
func makeTitle(md string) (title string) {
	// 获取title (第一个h1)
	titleBeg, titleEnd := strings.Index(md, "#")+1, strings.Index(md, "\n")
	if titleEnd == -1 { // 就一行, 特判
		titleEnd = len(md)
	}
	if titleBeg != -1 { // 没井号没标题
		title = strings.TrimSpace(md[titleBeg:titleEnd])
	}
	return
}

// 渲染的主要函数. 总体需要处理三个额外内容:
// 1. 渲染要吃一个斜杠\, 而latex还要用, 所以要重复1次
//      使用更友好的反斜杠(md内只需写一个, 就能保持转义两次).
//      markdown本身的特殊转义(#*_)除外.
// 2. 重复是全局的(普通段落里也可能有latex), 这导致一个问题:
//      pre代码中的斜杠\也被重复了两次!
//      所以需要一个未重复的版本, 渲染后进行还原
// 3. 渲染会将所有非html安全字符escape, 这在大部分地方都没问题, 也是必须做的.
//      但graph不支持 &lt; &gt; 这样的转义字符作为标识.
//      而且latex和graph的特殊功能段应该完整保留, 不要进行渲染
func makeHtml(md string) (htmlres string) {
	// ------------------------- 渲染准备 -----------------------------
	// 优化反斜杠\, 非转义#*_三个字符的反斜杠都一个变两个(或两个变四个)
	// 生成md2, 用于正式渲染. 而md保留给代码段进行安全渲染
	reg, _ := regexp.Compile(`\\([^#*_])`)
	md2 := strings.ReplaceAll(reg.ReplaceAllString(md, `\\$1`), `\\\`, `\\\\`)
	// ------------------------- 渲染进行 -----------------------------
	// 放入markdown渲染器
	// 一个通行版本(处理反斜杠)
	// 一个代码安全的版本(不处理反斜杠)
	var ext = blackfriday.WithExtensions(0 |
		blackfriday.Autolink |
		blackfriday.CommonExtensions |
		blackfriday.DefinitionLists |
		blackfriday.FencedCode |
		blackfriday.Footnotes |
		blackfriday.Tables,
	)
	htmlres = string(blackfriday.Run([]byte(md2), ext))
	htmlresSafe := string(blackfriday.Run([]byte(md), ext))
	// ------------------------- 渲染后处理 -----------------------------
	// 经历过了 BlackFriday 库的渲染
	// 现在特殊字符已经逃离, 代码段也已经不安全
	// 我们需要将:
	//   latex/graph还原为md初始样式(未进行html逃离的版本)
	//   pre还原为代码安全的版本
	// 注意这里的regGroup和htmlGroup是一个一维对二维的关系
	regGroup, htmlGroup := makeRestoreGroup(map[string]*string{
		`latex`: &md,
		`graph`: &md,
		`pre`:   &htmlresSafe,
	})
	// 用i遍历替换项(正则表达式)
	// 用j为每个替换项中的每个匹配段进行对应文本的替换
	for i := 0; i < len(regGroup); i++ {
		j := 0
		htmlres = regGroup[i].ReplaceAllStringFunc(htmlres, func(string) string {
			res := htmlGroup[i][j]
			j++
			return res
		})
	}
	return
}

// 一个封装, 将 正则表达式 和 对应匹配段目标内容组 联系起来.
func makeRestoreGroup(matches map[string]*string) (regGroup []*regexp.Regexp, htmlGroup [][]string) {
	for expr, findIn := range matches {
		reg := regexp.MustCompile(strings.ReplaceAll(`<PLACEHOLDER[\S\s]+?</PLACEHOLDER>`, `PLACEHOLDER`, expr))
		regGroup = append(regGroup, reg)
		htmlGroup = append(htmlGroup, reg.FindAllString(*findIn, -1))
	}
	return
}
