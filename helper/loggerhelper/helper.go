package loggerhelper

import "regexp"

func EscapeComma(str string) string {
	return regexp.MustCompile("\\s*,\\s*").ReplaceAllString(str, "|")
}
