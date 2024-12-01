package tmpbbs

import (
	"fmt"
	"html/template"
	"maps"
	"math"
	"slices"
	"strings"

	"time"

	"golang.org/x/text/message"
)

type emojiParser func(string) string

type displayPost struct {
	*post
	printer     *message.Printer
	emojiParser emojiParser
}

func newDisplayPost(post *post, printer *message.Printer, emojiParser emojiParser) *displayPost {
	return &displayPost{
		post:        post,
		printer:     printer,
		emojiParser: emojiParser,
	}
}

func (dp displayPost) BodyHTML() template.HTML {
	return template.HTML(markdownToHTML([]byte(dp.expandEmoji(dp.Body))))
}

func (dp displayPost) DisplayAuthor() string {
	return dp.expandEmoji(dp.Author)
}

func (dp displayPost) DisplayTitle() string {
	if dp.Title == "" {
		return fmt.Sprintf("#%d", dp.id)
	}

	return dp.expandEmoji(dp.Title)
}

func (dp displayPost) HasRepliesPage(page int, perPage int) bool {
	return page > 0 && page <= dp.repliesLastPage(perPage)
}

func (dp displayPost) NumRepliesLocalized() string {
	return dp.printer.Sprintf("%d replies", len(dp.Replies))
}

func (dp displayPost) ParentDisplayPost() *displayPost {
	return newDisplayPost(dp.Parent, dp.printer, dp.emojiParser)
}

func (dp displayPost) RepliesPage(page int, perPage int) []*displayPost {
	start := min((max(0, page-1))*perPage, len(dp.Replies))
	end := min(start+perPage, len(dp.Replies))

	var result []*displayPost
	for _, reply := range dp.Replies[start:end] {
		result = append(result, newDisplayPost(reply, dp.printer, dp.emojiParser))
	}

	return result
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
		if show[page] {
			links[i] = dp.printer.Sprintf(`<a href="%s">page %d</a>`, dp.repliesPageURL(page, "replies-start"), page)
		} else {
			links[i] = dp.printer.Sprintf("page %d", page)
		}
	}

	return template.HTML(fmt.Sprintf(`<li class="%s">%s</li>`, liClass, strings.Join(links, " / ")))
}

func (dp displayPost) TimeAgo() string {
	age := time.Since(dp.time)
	if age.Hours() < 1 {
		return dp.printer.Sprintf("%dm ago", int64(math.Round(age.Minutes())))
	}
	if age.Hours() >= 24 {
		return dp.printer.Sprintf("%dd ago", int64(math.Round(age.Hours()/24)))
	} else {
		return dp.printer.Sprintf("%dh ago", int64(math.Round(age.Hours())))
	}
}

func (dp displayPost) URL() string {
	return fmt.Sprintf("/%d", dp.id)
}

func (dp displayPost) expandEmoji(s string) string {
	if dp.emojiParser == nil {
		return s
	}

	return dp.emojiParser(s)
}

func (p post) repliesPageURL(page int, anchor string) string {
	if anchor != "" {
		anchor = fmt.Sprintf("#%s", anchor)
	}
	return fmt.Sprintf("/%d?p=%d%s", p.id, page, anchor)
}

func (p post) repliesLastPage(perPage int) int {
	return max(1, int(math.Ceil(float64(len(p.Replies))/float64(perPage))))
}

func (dp displayPost) repliesPageEndURL(perPage int, anchor string) string {
	return dp.repliesPageURL(dp.repliesLastPage(perPage), anchor)
}
