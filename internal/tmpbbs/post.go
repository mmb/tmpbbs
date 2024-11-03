package tmpbbs

import (
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

func (p post) IsOriginalPoster() bool {
	return p.Parent != nil && p.Parent.TripCode != "" && p.TripCode == p.Parent.TripCode
}
