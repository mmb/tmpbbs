package tmpbbs

import (
	"fmt"
	"html/template"
	"math"
	"time"

	"golang.org/x/text/message"
)

type displayPost struct {
	*post
	printer *message.Printer
}

func newDisplayPost(post *post, printer *message.Printer) *displayPost {
	return &displayPost{
		post:    post,
		printer: printer,
	}
}

func (dp displayPost) BodyHTML() template.HTML {
	return template.HTML(markdownToHTML([]byte(dp.Body)))
}

func (dp displayPost) DisplayTitle() string {
	if dp.Title == "" {
		return fmt.Sprintf("#%d", dp.id)
	}

	return dp.Title
}

func (dp displayPost) NumRepliesLocalized() string {
	if len(dp.Replies) == 1 {
		return dp.printer.Sprintf("%d reply", 1)
	} else {
		return dp.printer.Sprintf("%d replies", len(dp.Replies))
	}
}

func (dp displayPost) ParentDisplayPost() *displayPost {
	return newDisplayPost(dp.Parent, dp.printer)
}

func (dp displayPost) RepliesPage(page int, perPage int) []*displayPost {
	start := min((max(0, page-1))*perPage, len(dp.Replies))
	end := min(start+perPage, len(dp.Replies))

	var result []*displayPost
	for _, reply := range dp.Replies[start:end] {
		result = append(result, newDisplayPost(reply, dp.printer))
	}

	return result
}

func (dp displayPost) RepliesPageBeginURL() string {
	return dp.repliesPageURL(1)
}

func (dp displayPost) RepliesPagePrevURL(page int) string {
	return dp.repliesPageURL(max(page-1, 1))
}

func (dp displayPost) RepliesPageNextURL(page int, perPage int) string {
	return dp.repliesPageURL(min(dp.repliesLastPage(perPage), page+1))
}

func (dp displayPost) RepliesPageEndURL(perPage int) string {
	return dp.repliesPageURL(dp.repliesLastPage(perPage))
}

func (dp displayPost) TimeAgo() string {
	age := time.Since(dp.time)
	if age.Hours() < 1 {
		return fmt.Sprintf("%dm ago", int64(math.Round(age.Minutes())))
	}
	if age.Hours() >= 24 {
		return fmt.Sprintf("%dd ago", int64(math.Round(age.Hours()/24)))
	} else {
		return fmt.Sprintf("%dh ago", int64(math.Round(age.Hours())))
	}
}

func (dp displayPost) URL() string {
	return fmt.Sprintf("/%d", dp.id)
}

func (p post) repliesPageURL(page int) string {
	return fmt.Sprintf("/%d?p=%d", p.id, page)
}

func (p post) repliesLastPage(perPage int) int {
	return max(1, int(math.Ceil(float64(len(p.Replies))/float64(perPage))))
}
