package tmpbbs

import (
	"fmt"
	"math"
	"time"

	"github.com/google/uuid"
)

type post struct {
	time     time.Time
	Parent   *post
	Title    string
	Author   string
	Tripcode string
	Body     string
	uuid     string
	Replies  []*post
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
		uuid:     uuid.New().String(),
	}
}

func (p *post) IsOriginalPoster() bool {
	return p.Parent != nil && p.Parent.Tripcode != "" && p.Tripcode == p.Parent.Tripcode
}

func (p *post) URL() string {
	return "/p/" + p.uuid
}

func (p *post) repliesPageURL(page int, anchor string) string {
	if anchor != "" {
		anchor = "#" + anchor
	}

	return fmt.Sprintf("%s?p=%d%s", p.URL(), page, anchor)
}

func (p *post) hasRepliesPage(page int, perPage int) bool {
	return page > 0 && page <= p.repliesLastPage(perPage)
}

func (p *post) repliesLastPage(perPage int) int {
	return max(1, int(math.Ceil(float64(len(p.Replies))/float64(perPage))))
}
