package tmpbbs

import (
	"regexp"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
)

type markdownParser struct {
	blackfridayExtensions *blackfriday.Option
	blackfridayRenderer   *blackfriday.Option
	bluemondayPolicy      *bluemonday.Policy
}

func newMarkdownParser() *markdownParser {
	blackfridayExtensions := blackfriday.WithExtensions(
		blackfriday.Autolink |
			blackfriday.DefinitionLists |
			blackfriday.FencedCode |
			blackfriday.NoIntraEmphasis |
			blackfriday.Strikethrough |
			blackfriday.Tables)

	blackfridayRenderer := blackfriday.WithRenderer(
		blackfriday.NewHTMLRenderer(
			blackfriday.HTMLRendererParameters{
				// disable XHTML
				Flags: blackfriday.CommonHTMLFlags &^ blackfriday.UseXHTML,
			}))

	bluemondayPolicy := bluemonday.UGCPolicy()
	bluemondayPolicy.RequireNoReferrerOnLinks(true)
	bluemondayPolicy.AllowAttrs("class").Matching(regexp.MustCompile("^language-[a-zA-Z0-9]+$")).OnElements("code")

	return &markdownParser{
		blackfridayExtensions: &blackfridayExtensions,
		blackfridayRenderer:   &blackfridayRenderer,
		bluemondayPolicy:      bluemondayPolicy,
	}
}

func (mp markdownParser) parse(input string) string {
	unsafe := blackfriday.Run([]byte(input), *mp.blackfridayExtensions, *mp.blackfridayRenderer)

	return string(mp.bluemondayPolicy.SanitizeBytes(unsafe))
}
