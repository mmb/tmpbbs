package tmpbbs

import (
	"fmt"
	"html/template"
	"maps"
	"math"
	"slices"
	"strings"
	"time"

	"github.com/microcosm-cc/bluemonday"
	"golang.org/x/text/message"
)

var bluemondayStrictPolicy = bluemonday.StrictPolicy()

type parser interface {
	parse(input string) string
}

type displayPost struct {
	*post
	Printer             *message.Printer
	basicEmojiParser    parser
	wrappingEmojiParser parser
	markdownParser      parser
}

func newDisplayPost(post *post, printer *message.Printer, basicEmojiParser parser, wrappingEmojiParser parser, markdownParser parser) *displayPost {
	return &displayPost{
		post:                post,
		Printer:             printer,
		basicEmojiParser:    basicEmojiParser,
		wrappingEmojiParser: wrappingEmojiParser,
		markdownParser:      markdownParser,
	}
}

func (dp displayPost) BodyHTML() template.HTML {
	return template.HTML(dp.expandEmoji(dp.markdownParser.parse(dp.Body), dp.wrappingEmojiParser))
}

func (dp displayPost) DisplayAuthor() template.HTML {
	return template.HTML(dp.sanitizeAndExpandEmoji(dp.Author, dp.wrappingEmojiParser))
}

func (dp displayPost) DisplayTitle() template.HTML {
	title := dp.sanitizeAndExpandEmoji(dp.Title, dp.wrappingEmojiParser)

	if title == "" {
		return template.HTML(dp.emptyTitle())
	}

	return template.HTML(title)
}

func (dp displayPost) NumReplies() string {
	return dp.Printer.Sprintf("%d replies", len(dp.Replies))
}

func (dp displayPost) PageTitle() string {
	title := dp.sanitizeAndExpandEmoji(dp.Title, dp.basicEmojiParser)

	if title == "" {
		return dp.emptyTitle()
	}

	return title
}

func (dp displayPost) ParentDisplayPost() *displayPost {
	return newDisplayPost(dp.Parent, dp.Printer, dp.basicEmojiParser, dp.wrappingEmojiParser, dp.markdownParser)
}

func (dp displayPost) RepliesNav(currentPage int, perPage int, liClass string) template.HTML {
	if len(dp.Replies) == 0 {
		return ""
	}

	firstPage := 1
	prevPage := max(currentPage-1, 1)
	nextPage := min(dp.repliesLastPage(perPage), currentPage+1)
	lastPage := dp.repliesLastPage(perPage)

	show := make(map[int]bool, 5)
	show[firstPage] = true
	show[prevPage] = true
	show[nextPage] = true
	show[lastPage] = true
	show[currentPage] = false

	if len(show) == 1 {
		return ""
	}

	pages := slices.Sorted(maps.Keys(show))

	links := make([]string, len(pages))
	for i, page := range pages {
		pageText := dp.Printer.Sprintf("page %d", page)
		if show[page] {
			links[i] = fmt.Sprintf(`<a href="%s">%s</a>`, dp.repliesPageURL(page, "replies-start"), pageText)
		} else {
			links[i] = pageText
		}
	}

	return template.HTML(fmt.Sprintf(`<li class="%s">%s</li>`, liClass, strings.Join(links, " / ")))
}

func (dp displayPost) RepliesPage(page int, perPage int) []*displayPost {
	start := min((max(0, page-1))*perPage, len(dp.Replies))
	end := min(start+perPage, len(dp.Replies))

	result := make([]*displayPost, end-start)
	for i, reply := range dp.Replies[start:end] {
		result[i] = newDisplayPost(reply, dp.Printer, dp.basicEmojiParser, dp.wrappingEmojiParser, dp.markdownParser)
	}

	return result
}

func (dp displayPost) TimeAgo() string {
	age := time.Since(dp.time)
	if age.Hours() < 1 {
		return dp.Printer.Sprintf("%dm ago", int64(math.Round(age.Minutes())))
	}
	if age.Hours() >= 24 {
		return dp.Printer.Sprintf("%dd ago", int64(math.Round(age.Hours()/24)))
	} else {
		return dp.Printer.Sprintf("%dh ago", int64(math.Round(age.Hours())))
	}
}

func (dp displayPost) emptyTitle() string {
	return fmt.Sprintf("#%d", dp.id)
}

func (dp displayPost) expandEmoji(input string, parser parser) string {
	if parser == nil {
		return input
	}

	return parser.parse(input)
}

func (dp displayPost) sanitizeAndExpandEmoji(input string, parser parser) string {
	return dp.expandEmoji(strings.TrimSpace(string(bluemondayStrictPolicy.SanitizeBytes([]byte(input)))), parser)
}
