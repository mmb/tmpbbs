package tmpbbs

import (
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
)

var blackfridayExtensions = blackfriday.WithExtensions(
	blackfriday.Autolink |
		blackfriday.HardLineBreak)

var bluemondayPolicy = bluemonday.UGCPolicy()

func init() {
	bluemondayPolicy.RequireNoReferrerOnLinks(true)
}

func markdownToHTML(markdown []byte) []byte {
	unsafe := blackfriday.Run(markdown, blackfridayExtensions)

	return bluemondayPolicy.SanitizeBytes(unsafe)
}
