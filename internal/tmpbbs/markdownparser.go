package tmpbbs

import (
	"bytes"
	"regexp"

	"github.com/microcosm-cc/bluemonday"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

type markdownParser struct {
	markdown         goldmark.Markdown
	bluemondayPolicy *bluemonday.Policy
}

var allowedCodeClass = regexp.MustCompile("^language-[a-zA-Z0-9]+$")

func newMarkdownParser() *markdownParser {
	bluemondayPolicy := bluemonday.UGCPolicy()
	bluemondayPolicy.RequireNoReferrerOnLinks(true)
	bluemondayPolicy.AllowAttrs("class").Matching(allowedCodeClass).OnElements("code")

	return &markdownParser{
		markdown:         goldmark.New(goldmark.WithExtensions(extension.DefinitionList, extension.GFM)),
		bluemondayPolicy: bluemondayPolicy,
	}
}

func (mp *markdownParser) parse(input string) string {
	var unsafe bytes.Buffer
	mp.markdown.Convert([]byte(input), &unsafe) //nolint:errcheck,gosec // bytes.Buffer.Write() always return nil for error

	return string(mp.bluemondayPolicy.SanitizeBytes(unsafe.Bytes()))
}
