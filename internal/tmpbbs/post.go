package tmpbbs

import (
	"fmt"
	"math"
	"time"
)

type post struct {
	time     time.Time
	Parent   *post
	Title    string
	Author   string
	Tripcode string
	Body     string
	Replies  []*post
	id       int
}

func newPost(title string, author string, body string, tripcoder *Tripcoder) *post {
	var tripcode string
	if tripcoder != nil {
		author, tripcode = tripcoder.code(author)
	}

	return &post{
		Title:    title,
		Author:   author,
		Tripcode: tripcode,
		Body:     body,
		time:     time.Now(),
	}
}

func (p *post) IsOriginalPoster() bool {
	return p.Parent != nil && p.Parent.Tripcode != "" && p.Tripcode == p.Parent.Tripcode
}

func (p *post) URL() string {
	return fmt.Sprintf("/%d", p.id)
}

func (p *post) repliesPageURL(page int, anchor string) string {
	if anchor != "" {
		anchor = "#" + anchor
	}

	return fmt.Sprintf("/%d?p=%d%s", p.id, page, anchor)
}

func (p *post) hasRepliesPage(page int, perPage int) bool {
	return page > 0 && page <= p.repliesLastPage(perPage)
}

func (p *post) repliesLastPage(perPage int) int {
	return max(1, int(math.Ceil(float64(len(p.Replies))/float64(perPage))))
}
