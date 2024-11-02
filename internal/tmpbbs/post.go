package tmpbbs

import (
	"fmt"
	"html/template"
	"math"
	"time"
)

type post struct {
	id       int
	Title    string
	Author   string
	TripCode string
	Body     string
	Parent   *post
	Replies  []*post
	time     time.Time
}

func newPost(title string, author string, body string, tripCoder *tripCoder) *post {
	var tripCode string
	if tripCoder != nil {
		author, tripCode = tripCoder.code(author)
	}
	return &post{
		Title:    title,
		Author:   author,
		TripCode: tripCode,
		Body:     body,
		time:     time.Now(),
	}
}

func (p post) BodyHTML() template.HTML {
	return template.HTML(markdownToHTML([]byte(p.Body)))
}

func (p post) TimeAgo() string {
	age := time.Since(p.time)
	if age.Hours() < 1 {
		return fmt.Sprintf("%dm ago", int64(math.Round(age.Minutes())))
	}
	if age.Hours() >= 24 {
		return fmt.Sprintf("%dd ago", int64(math.Round(age.Hours()/24)))
	} else {
		return fmt.Sprintf("%dh ago", int64(math.Round(age.Hours())))
	}
}

func (p post) DisplayTitle() string {
	if p.Title == "" {
		return fmt.Sprintf("#%d", p.id)
	}

	return p.Title
}

func (p post) URL() string {
	return fmt.Sprintf("/%d", p.id)
}

func (p post) BeginRepliesPageURL() string {
	return p.repliesPageURL(1)
}

func (p post) PrevRepliesPageURL(page int) string {
	return p.repliesPageURL(max(page-1, 1))
}

func (p post) NextRepliesPageURL(page int, perPage int) string {
	return p.repliesPageURL(min(p.repliesLastPage(perPage), page+1))
}

func (p post) EndRepliesPageURL(perPage int) string {
	return p.repliesPageURL(p.repliesLastPage(perPage))
}

func (p post) RepliesPage(page int, perPage int) []*post {
	start := min((max(0, page-1))*perPage, len(p.Replies))
	end := min(start+perPage, len(p.Replies))

	return p.Replies[start:end]
}

func (p post) repliesPageURL(page int) string {
	return fmt.Sprintf("/%d?p=%d", p.id, page)
}

func (p post) repliesLastPage(perPage int) int {
	return max(1, int(math.Ceil(float64(len(p.Replies))/float64(perPage))))
}
