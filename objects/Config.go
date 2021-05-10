package objects

import (
	"html/template"
	"os"
	"regexp"
	"runtime"

	"github.com/rickonono3/m2obj"
	"github.com/rickonono3/m2obj/m2json"
)

var CWD string
var Config *m2obj.Object
var ConfigFile = "config.json"
var ConfigFileSyncer *m2obj.FileSyncer
var DefaultConfig = m2obj.New(m2obj.Group{
	"Version": "0.0.0",
	"Blog": m2obj.Group{
		"CDN":     "/",
		"BGImg":   "img/bg.jpg",
		"IconImg": "img/icon.svg",
		"Favicon": "favicon.ico",
		"Icons": m2obj.Group{
			"Dir":     "&#xe806;",
			"Article": "&#xe809;",
			"File":    "&#xe7f5;",
		},
	},
	"Info": m2obj.Group{
		"Email":    "rickonono3@rickonono3.top",
		"QQ":       "1070101572",
		"Telegram": "RickoNoNo3",
	},
	"AdminPSWD":  "",
	"ServerPort": 13808,
	"IsInDebug":  false,
	"LogFile": m2obj.Group{
		"ConsoleLog": "",
		"WebLog":     "web_log.csv",
	},
})

func makeCWD() (cwd string) {
	if blogRoot := os.Getenv("BlogRoot"); blogRoot != "" {
		cwd = regexp.MustCompile("/?$").ReplaceAllString(blogRoot, "/")
	} else {
		cwd = os.Args[0]
		if runtime.GOOS == "windows" {
			cwd = regexp.MustCompile(`^(.*[\\])[^\\]*$`).ReplaceAllString(cwd, "$1")
		} else {
			cwd = regexp.MustCompile(`^(.*[/])[^/]*$`).ReplaceAllString(cwd, "$1")
		}
		if cwd == os.Args[0] {
			return "./"
		}
	}
	return
}

func initConfig() (err error) {
	CWD = makeCWD()
	Config = DefaultConfig.Clone()
	ConfigFileSyncer = m2obj.NewFileSyncer(CWD+ConfigFile, m2json.Formatter{})
	ConfigFileSyncer.BindObject(Config)
	ConfigFileSyncer.Load()
	Config.MustGet("Blog.Icons").GroupForeach(func(key string, obj *m2obj.Object) error {
		obj.SetVal(template.HTML(obj.ValStr()))
		return nil
	})
	return
}
