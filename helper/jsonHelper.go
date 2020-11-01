package helper

import (
	"bytes"
	"encoding/json"
)

func GenerateJson(v interface{}) (jsonRes string) {
	buf := bytes.Buffer{}
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	encoder.Encode(v)
	return buf.String()
}
