package tmpbbs

import (
	"fmt"
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

func newPost(title string, author string, body string, tripCoder *TripCoder) *post {
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

func (p post) IsOriginalPoster() bool {
	return p.Parent != nil && p.Parent.TripCode != "" && p.TripCode == p.Parent.TripCode
}

func (p post) URL() string {
	return fmt.Sprintf("/%d", p.id)
}

func (p post) repliesPageURL(page int, anchor string) string {
	if anchor != "" {
		anchor = "#" + anchor
	}

	return fmt.Sprintf("/%d?p=%d%s", p.id, page, anchor)
}

func (p post) hasRepliesPage(page int, perPage int) bool {
	return page > 0 && page <= p.repliesLastPage(perPage)
}

func (p post) repliesLastPage(perPage int) int {
	return max(1, int(math.Ceil(float64(len(p.Replies))/float64(perPage))))
}
