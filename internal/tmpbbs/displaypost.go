package tmpbbs

import (
	"fmt"
	"html/template"
	"math"
	"sort"
	"strings"

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

func (dp displayPost) RepliesNav(currentPage int, perPage int) template.HTML {
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
		return template.HTML(dp.NumRepliesLocalized())
	}

	var pages []int
	for page := range show {
		pages = append(pages, page)
	}
	sort.Ints(pages)

	var links []string
	for _, page := range pages {
		if show[page] {
			links = append(links, dp.printer.Sprintf(`<a href="%s">page %d</a>`, dp.repliesPageURL(page), page))
		} else {
			links = append(links, dp.printer.Sprintf("page %d", page))
		}
	}

	return template.HTML(fmt.Sprintf("%s, %s", dp.NumRepliesLocalized(), strings.Join(links, " / ")))
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

func (p post) repliesPageURL(page int) string {
	return fmt.Sprintf("/%d?p=%d", p.id, page)
}

func (p post) repliesLastPage(perPage int) int {
	return max(1, int(math.Ceil(float64(len(p.Replies))/float64(perPage))))
}

func (dp displayPost) repliesPageEndURL(perPage int) string {
	return dp.repliesPageURL(dp.repliesLastPage(perPage))
}
