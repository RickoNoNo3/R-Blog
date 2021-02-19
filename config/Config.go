/**
 * Debug模式:
 * 1. 模板引擎在每次访问时刷新
 */
package config

import (
	"os"
	"regexp"
	"runtime"

	"rickonono3/r-blog/mytype"
)

var Config *mytype.Object

func init() {
	cwd := os.Args[0]
	cwd = "./"
	if runtime.GOOS == "windows" {
		cwd = regexp.MustCompile(`^(.*[\\])[^\\]*$`).ReplaceAllString(cwd, "$1")
	} else {
		cwd = regexp.MustCompile(`^(.*[/])[^/]*$`).ReplaceAllString(cwd, "$1")
	}
	Config = mytype.NewGroup(mytype.ObjectList{
		"Cwd": mytype.NewValue(cwd),
		"Blog": mytype.NewGroup(mytype.ObjectList{
			// "CDN":   mytype.NewValue("http://81.70.142.172/"),
			"CDN":   mytype.NewValue("http://localhost:13808/"),
			"BGImg": mytype.NewValue("img/bg.jpg"),
			"Icons": mytype.NewGroup(mytype.ObjectList{
				"Dir":     mytype.NewValue("folder"),
				"Article": mytype.NewValue("description"),
				"File":    mytype.NewValue("file"),
			}),
		}),
		"AdminPSWD": mytype.NewValue("tydhc199963+"),
		"IsInDebug": mytype.NewValue(true),
	})
}

func Get(key string) *mytype.Object {
	return Config.Get(key)
}

func Set(key string, val *mytype.Object) *mytype.Object {
	return Config.Set(key, val)
}
