package drawer

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"strings"
	"text/template"
	"time"

	"github.com/chromedp/chromedp"
)

// Render 传入一个Parse后的字符串, 这时应该只做了md解析, 是一个html片段,
// 本函数会用真实的浏览器环境, 加载各类JS引擎(highlight.js/mathjax/viz.js).
func Render(htmlori string) (htmlres string) {
	noRender := true
	defer func() {
		if noRender {
			htmlres = "<blockquote><p>警告: 此页面未经高级渲染</p></blockquote>" + htmlres
		}
	}()
	htmlres = htmlori
	// 模板填充
	// 先往模板页(drawer.html)里填充传入的半成品html片段(htmlori).
	tmpl, err := template.ParseFiles("drawer/web/drawer.html")
	if err != nil {
		fmt.Println(err)
		return
	}
	buf := new(bytes.Buffer)
	err = tmpl.Execute(buf, htmlori)
	if err != nil {
		fmt.Println(err)
		return
	}
	htmlexec := buf.String()

	// 创建一个临时服务器(对就是普通web服务器),
	// 把这个模板页, 还有配套的资源都放进去.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" { // 模板页面
			io.Copy(w, bytes.NewBufferString(htmlexec))
		} else { // 依赖文件
			file, err := os.Open("drawer/web" + r.URL.Path)
			if err != nil {
				fmt.Println(err)
				return
			}
			defer file.Close()
			io.Copy(w, file)
		}
	}))
	defer ts.Close()

	// 再用 headless 的 chrome 执行渲染, 和页面上的js配合, 然后保存渲染结果.
	preCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ctx, cancel := chromedp.NewContext(preCtx)
	defer cancel()
	err = chromedp.Run(ctx,
		chromedp.Navigate(ts.URL),       // 导航到模板页面
		chromedp.WaitReady("body>done"), // 页面里js会产生一个done标签, 代表渲染结束
		chromedp.InnerHTML("body", &htmlres),
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	// 删除 done 标签,
	// 删除带有 draw="not-need" 属性的 script 标签.
	reg, _ := regexp.Compile(`<done></done>|<script[^>]*draw="not-need"[\S\s]+?</script>`)
	htmlres = reg.ReplaceAllString(htmlres, "")
	// 连续多个空行压缩成一行
	reg, _ = regexp.Compile(`\n{2,}`)
	htmlres = reg.ReplaceAllString(htmlres, "\n")
	// Trim
	htmlres = strings.Trim(htmlres, "\n 　\t")
	noRender = false
	return
}
