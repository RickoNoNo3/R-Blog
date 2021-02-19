package server

import (
	"regexp"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/html"
	"github.com/tdewolff/minify/v2/json"
	"github.com/tdewolff/minify/v2/svg"
	"github.com/tdewolff/minify/v2/xml"
)

var M = minify.New()

func loadMinify() {
	M.Add("text/html", &html.Minifier{
		KeepDocumentTags: true,
		KeepEndTags:      true,
		KeepQuotes:       true,
		KeepWhitespace:   true,
	})
	M.AddFunc("image/svg+xml", svg.Minify)
	M.AddFuncRegexp(regexp.MustCompile(`[/+]json$`), json.Minify)
	M.AddFuncRegexp(regexp.MustCompile(`[/+]xml$`), xml.Minify)
}
