/**
 * Debug模式:
 * 1. 模板引擎在每次访问时刷新
 */
package objects

import (
	"os"
	"regexp"
	"runtime"

	"rickonono3/r-blog/mytype"
)

var Config *mytype.Object

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
