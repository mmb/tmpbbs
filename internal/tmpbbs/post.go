package tmpbbs

import (
	"fmt"
	"math"
	"time"
)

type post struct {
	id      int
	Title   string
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
		Title:  title,
		Author: author,
		Body:   body,
		time:   time.Now(),
	}
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

func (p post) URL() string {
	return fmt.Sprintf("/%d", p.id)
}
