package tmpbbs

import (
	"regexp"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
)

var (
	blackfridayExtensions = blackfriday.WithExtensions(
		blackfriday.Autolink |
			blackfriday.DefinitionLists |
			blackfriday.FencedCode |
			blackfriday.NoIntraEmphasis |
			blackfriday.Strikethrough |
			blackfriday.Tables)

	blackfridayRenderer = blackfriday.WithRenderer(
		blackfriday.NewHTMLRenderer(
			blackfriday.HTMLRendererParameters{
				// disable XHTML
				Flags: blackfriday.CommonHTMLFlags &^ blackfriday.UseXHTML,
			}))

	bluemondayPolicy = bluemonday.UGCPolicy()
)

func InitBluemonday() {
	bluemondayPolicy.RequireNoReferrerOnLinks(true)
	bluemondayPolicy.AllowAttrs("class").Matching(regexp.MustCompile("^language-[a-zA-Z0-9]+$")).OnElements("code")
}

type markdownParser struct{}

func newMarkdownParser() *markdownParser {
	return &markdownParser{}
}

func (mp markdownParser) parse(input string) string {
	unsafe := blackfriday.Run([]byte(input), blackfridayExtensions, blackfridayRenderer)

	return string(bluemondayPolicy.SanitizeBytes(unsafe))
}
