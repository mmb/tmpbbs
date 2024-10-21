package tmpbbs

import (
	"fmt"
	"html/template"
	"math"
	"time"
)

type post struct {
	id      int
	title   string
	Author  string
	Body    string
	Parent  *post
	Replies []*post
	time    time.Time
}

func newPost(title string, author string, body string, tripCoder *tripCoder) *post {
	if tripCoder != nil {
		author = tripCoder.code(author)
	}
	return &post{
		title:  title,
		Author: author,
		Body:   body,
		time:   time.Now(),
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

func (p post) Title() string {
	if p.title != "" {
		return p.title
	}

	return fmt.Sprintf("#%d", p.id)
}

func (p post) URL() string {
	return fmt.Sprintf("/%d", p.id)
}
