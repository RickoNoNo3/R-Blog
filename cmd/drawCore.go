package cmd

import (
	"io/ioutil"
	"os"

	"../drawer"
)

// 用文件(或者缺省为stdin)作为文章来渲染
func DrawCore(inputFile *os.File) (out string) {
	input, err := ioutil.ReadAll(inputFile)
	if err == nil {
		_, out = drawer.Draw(string(input))
	}
	return
}
