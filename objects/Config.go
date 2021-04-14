/**
 * Debug模式:
 * 1. 模板引擎在每次访问时刷新
 */
package objects

import (
	"html/template"
	"os"
	"regexp"
	"runtime"

	"github.com/rickonono3/m2obj"
	"github.com/rickonono3/m2obj/m2json"
)

var Config *m2obj.Object
var ConfigFile = "config.json"
var DefaultConfig = m2obj.New(m2obj.Group{
	"Version": "0.0.0",
	"Cwd":     "./",
	"Blog": m2obj.Group{
		"CDN":     "/",
		"BGImg":   "img/bg.jpg",
		"IconImg": "img/icon.svg",
		"Icons": m2obj.Group{
			"Dir":     "&#xe806;",
			"Article": "&#xe809;",
			"File":    "&#xe7f5;",
		},
	},
	"AdminPSWD": "",
	"IsInDebug": false,
})

func makeCwd() (cwd string) {
	cwd = os.Args[0]
	if runtime.GOOS == "windows" {
		cwd = regexp.MustCompile(`^(.*[\\])[^\\]*$`).ReplaceAllString(cwd, "$1")
	} else {
		cwd = regexp.MustCompile(`^(.*[/])[^/]*$`).ReplaceAllString(cwd, "$1")
	}
	if cwd == os.Args[0] {
		return "./"
	}
	return
}

func initConfig() {
	var (
		err error
		cwd = makeCwd()
	)
	if blogRoot := os.Getenv("BlogRoot"); blogRoot != "" {
		cwd = regexp.MustCompile("/?$").ReplaceAllString(blogRoot, "/")
	}
	if Config, err = m2json.New().LoadFromFile(cwd + ConfigFile); err != nil {
		Config = DefaultConfig
		m2json.New().SaveToFile(Config, cwd+ConfigFile)
	}
	Config.Set("Cwd", cwd)
	{
		dir := Config.MustGet("Blog.Icons.Dir")
		article := Config.MustGet("Blog.Icons.Article")
		file := Config.MustGet("Blog.Icons.File")
		dir.SetVal(template.HTML(dir.ValStr()))
		article.SetVal(template.HTML(article.ValStr()))
		file.SetVal(template.HTML(file.ValStr()))
	}
}

func SaveConfig() {
	m2json.New().SaveToFile(Config, Config.MustGet("Cwd").ValStr()+ConfigFile)
}
