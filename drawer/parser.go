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

// 对md进行格式处理
func makeMd(md string) string {
	var reg *regexp.Regexp
	md = strings.ReplaceAll(md, "\r\n", "\n")
	// $$...$$ => <latex>\[...\]</latex>
	reg, _ = regexp.Compile(`(^|[^\\$])\${2}(([\n\r]|.)*?)([^\\$])\${2}`)
	md = reg.ReplaceAllString(md, "$1<latex>\\[$2$4\\]</latex>\n")
	// $...$ => <latex-inner>\(...\)</latex-inner>
	reg, _ = regexp.Compile(`(^|[^\\$])\$(([\n\r]|.)*?)([^\\$])\$`)
	md = reg.ReplaceAllString(md, "$1<latex-inner>\\($2$4\\)</latex-inner>")
	// trim
	md = strings.Trim(md, "\n\r 　\t")
	// fmt.Println(md)
	return md
}

// 拿标题
func makeTitle(md string) (title string) {
	// 获取title (第一个h1)
	firstLine := strings.TrimSpace(strings.Split(md, "\n")[0])
	titleBeg, titleEnd := strings.Index(firstLine, "#"), strings.Index(md, "\n")
	if titleEnd == -1 { // 就一行, 特判
		titleEnd = len(firstLine)
	}
	if titleBeg != -1 { // 没井号没标题
		return strings.TrimSpace(firstLine[titleBeg+1 : titleEnd])
	} else {
		return firstLine
	}
}

// 渲染的主要函数. 总体需要处理三个额外内容:
//
// 1. 渲染要吃一个斜杠\, 而latex还要用, 所以要重复1次
//      使用更友好的反斜杠(md内只需写一个, 就能保持转义两次).
//      markdown本身的特殊转义(#*$_)除外.
// 2. 重复是全局的(普通段落里也可能有latex), 这导致一个问题:
//      pre代码中的斜杠\也被重复了两次!
//      所以需要一个未重复的版本, 渲染后进行还原
// 3. 渲染会将所有非html安全字符escape, 这在大部分地方都没问题, 也是必须做的.
//      但graph不支持 &lt; &gt; 这样的转义字符作为标识.
//      而且latex和graph的特殊功能段应该完整保留, 不要进行渲染
func makeHtml(md string) (htmlres string) {
	// ------------------------- 渲染准备 -----------------------------
	// 优化反斜杠\, 一个变两个(或两个变四个), 生成md2
	// 先用md2进行基础渲染(暴力过解析器), 再用md中的原文来修正md2和暴力过解析器带来的问题
	reg, _ := regexp.Compile(`\\([^#*$_])`)
	md2 := reg.ReplaceAllString(md, `\\$1`)
	md2 = strings.ReplaceAll(md2, `\\\`, `\\\\`)
	md2 = strings.ReplaceAll(md2, "&gt;", ">")
	md2 = strings.ReplaceAll(md2, "&lt;", "<")
	// ------------------------- 渲染进行 -----------------------------
	// 放入markdown渲染器
	// 一个通行版本(处理反斜杠)
	// 一个代码安全的版本(不处理反斜杠)
	var ext = blackfriday.WithExtensions(0 |
		blackfriday.Autolink |
		blackfriday.DefinitionLists |
		blackfriday.FencedCode |
		blackfriday.Footnotes |
		blackfriday.Tables |
		blackfriday.AutoHeadingIDs |
		blackfriday.HeadingIDs |
		blackfriday.LaxHTMLBlocks |
		blackfriday.Strikethrough |
		blackfriday.NoIntraEmphasis,
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
		`latex`:  &md,
		`latex-inner`: &md,
		`graph`:  &md,
		`pre`:    &htmlresSafe,
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
	// 把所有没嵌套p的dd都嵌套一个p (因为复杂定义中dd会自动渲染出p, 保持一致性)
	htmlres = regexp.MustCompile(`<dd>([^<][^p][^>])`).ReplaceAllString(htmlres, `<dd><p>$1`)
	htmlres = regexp.MustCompile(`([^<][^/][^p][^>])</dd>`).ReplaceAllString(htmlres, `$1</p></dd>`)
	return
}

// 一个封装, 将 正则表达式 和 对应匹配段目标内容组 联系起来.
func makeRestoreGroup(matches map[string]*string) (regGroup []*regexp.Regexp, htmlGroup [][]string) {
	for expr, findIn := range matches {
		reg := regexp.MustCompile(strings.ReplaceAll(`<PLACEHOLDER( +[^>]*)*?>([\n\r]|.)+?</PLACEHOLDER>`, `PLACEHOLDER`, expr))
		regGroup = append(regGroup, reg)
		htmlGroup = append(htmlGroup, reg.FindAllString(*findIn, -1))
	}
	return
}
