/*
Package tmpbbs provides all functionality for the tmpbbs binary.

It includes http and gRPC servers and handlers, HTML templating,
localization, and data store.
*/
package tmpbbs

import (
	"fmt"
	"html/template"
	"slices"
	"strings"
	"time"

	"github.com/microcosm-cc/bluemonday"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type parser interface {
	parse(input string) string
}

type displayPost struct {
	*post

	basicEmojiParser    parser
	Language            language.Tag
	markdownParser      parser
	Printer             *message.Printer
	wrappingEmojiParser parser
}

var strictPolicy = bluemonday.StrictPolicy() //nolint:gochecknoglobals // shared policy

func newDisplayPost(post *post, lang language.Tag, printer *message.Printer, basicEmojiParser parser,
	wrappingEmojiParser parser, markdownParser parser,
) *displayPost {
	return &displayPost{
		basicEmojiParser:    basicEmojiParser,
		Language:            lang,
		markdownParser:      markdownParser,
		post:                post,
		Printer:             printer,
		wrappingEmojiParser: wrappingEmojiParser,
	}
}

// BodyHTML returns a [template.HTML] of the post's body with Markdown and emoji
// expanded.
func (dp *displayPost) BodyHTML() template.HTML {
	return template.HTML(
		dp.expandEmoji(dp.markdownParser.parse(dp.Body), dp.wrappingEmojiParser)) // #nosec G203 -- sanitized
}

// DisplayAuthor returns a [template.HTML] of the posts's author with emoji
// expanded wrapped in a bdi tag and prefixed with "by ".
func (dp *displayPost) DisplayAuthor() template.HTML {
	authorSanitized := dp.sanitizeAndExpandEmoji(dp.Author, dp.wrappingEmojiParser)
	if authorSanitized == "" {
		return ""
	}

	author := `<bdi class="author">` + authorSanitized + "</bdi>"

	return template.HTML(dp.Printer.Sprintf("by %s", author)) // #nosec G203 -- sanitized
}

// DisplayTitle returns a [template.HTML] of the posts's title with emoji
// expanded.
func (dp *displayPost) DisplayTitle() template.HTML {
	title := dp.sanitizeAndExpandEmoji(dp.Title, dp.wrappingEmojiParser)

	if title == "" {
		title = dp.emptyTitle()
	}

	return template.HTML(title) // #nosec G203 -- either sanitized above or no user input if empty
}

// LanguageDirection returns the HTML text direction of the user's browser
// language ("ltr" for left-to-right or "rtl" for right-to-left). Unrecognized
// scripts return "ltr".
func (dp *displayPost) LanguageDirection() string {
	script, _ := dp.Language.Script()
	switch script.String() {
	case "Arab", "Hebr":
		return "rtl"
	default:
		return "ltr"
	}
}

// NumReplies returns how many replies the post has in human language.
func (dp *displayPost) NumReplies() string {
	if dp.Replies.Len() == 0 {
		return ""
	}

	return dp.Printer.Sprintf("(%d replies)", dp.Replies.Len())
}

// PageTitle returns a [template.HTML] containing the page title for the post's
// page.
func (dp *displayPost) PageTitle() template.HTML {
	title := dp.sanitizeAndExpandEmoji(dp.Title, dp.basicEmojiParser)

	if title == "" {
		title = dp.emptyTitle()
	}

	return template.HTML(title) // #nosec G203 -- either sanitized above or no user input if empty
}

// ParentDisplayPost returns the displayPost of the post's parent.
func (dp *displayPost) ParentDisplayPost() *displayPost {
	return newDisplayPost(dp.Parent, dp.Language, dp.Printer, dp.basicEmojiParser, dp.wrappingEmojiParser,
		dp.markdownParser)
}

// RepliesNav returns a [template.HTML] with navigation links based on the
// current page.
func (dp *displayPost) RepliesNav(currentPage int, perPage int, liClass string) template.HTML {
	if dp.Replies.Len() == 0 {
		return ""
	}

	firstPage := 1
	prevPage := max(currentPage-1, 1)
	lastPage := dp.repliesLastPage(perPage)
	nextPage := min(lastPage, currentPage+1)
	pages := []int{firstPage, prevPage, currentPage, nextPage, lastPage}
	slices.Sort(pages)

	pages = slices.Compact(pages)
	if len(pages) == 1 {
		return ""
	}

	links := make([]string, len(pages))
	for i, page := range pages {
		pageText := dp.Printer.Sprintf("page %d", page)
		if page == currentPage {
			links[i] = pageText
		} else {
			links[i] = fmt.Sprintf("<a href=%q>%s</a>", dp.repliesPageURL(page, "replies-start"), pageText)
		}
	}

	return template.HTML(
		fmt.Sprintf("<li class=%q>%s</li>", liClass, strings.Join(links, " / "))) // #nosec G203 -- no user input
}

// RepliesPage returns a page of the posts's replies.
func (dp *displayPost) RepliesPage(page int, perPage int) []*displayPost {
	start := min(max(0, page-1)*perPage, dp.Replies.Len())
	end := min(start+perPage, dp.Replies.Len())
	result := make([]*displayPost, end-start)

	current := dp.Replies.Front()
	for range start {
		current = current.Next()
	}

	for i := 0; i < perPage && current != nil; i++ {
		result[i] = newDisplayPost(current.Value.(*post), //nolint:errcheck,forcetypeassert // only one type
			dp.Language, dp.Printer, dp.basicEmojiParser, dp.wrappingEmojiParser, dp.markdownParser)
		current = current.Next()
	}

	return result
}

// TimeRFC3339 returns the post's time in UTC formatted as RFC3339.
func (dp *displayPost) TimeRFC3339() string {
	return dp.Time.UTC().Format(time.RFC3339)
}

func (dp *displayPost) emptyTitle() string {
	return strings.ToLower(fmt.Sprintf("%s-%s", dp.id[20:24], dp.id[24:26]))
}

func (dp *displayPost) expandEmoji(input string, parser parser) string {
	if parser == nil {
		return input
	}

	return parser.parse(input)
}

func (dp *displayPost) sanitizeAndExpandEmoji(input string, parser parser) string {
	return dp.expandEmoji(strings.TrimSpace(strictPolicy.Sanitize(input)), parser)
}
