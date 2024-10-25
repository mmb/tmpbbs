package tmpbbs

import (
	"regexp"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
)

var blackfridayExtensions = blackfriday.WithExtensions(
	blackfriday.Autolink |
		blackfriday.FencedCode |
		blackfriday.HardLineBreak)

var bluemondayPolicy = bluemonday.UGCPolicy()

func init() {
	bluemondayPolicy.RequireNoReferrerOnLinks(true)
	bluemondayPolicy.AllowAttrs("class").Matching(regexp.MustCompile("^language-[a-zA-Z0-9]+$")).OnElements("code")
}

func markdownToHTML(markdown []byte) []byte {
	unsafe := blackfriday.Run(markdown, blackfridayExtensions)

	return bluemondayPolicy.SanitizeBytes(unsafe)
}
