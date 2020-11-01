package helper

import "strings"

func CutToOneLine(str string) (res string) {
	return strings.TrimSpace(strings.ReplaceAll(strings.Split(str, "\n")[0], "\r", ""))
}