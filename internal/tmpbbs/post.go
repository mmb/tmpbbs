package tmpbbs

import (
	"fmt"
	"math"
	"time"

	"github.com/google/uuid"
)

type post struct {
	time      time.Time
	Parent    *post
	Title     string
	Author    string
	Tripcode  string
	Body      string
	uuid      string
	Replies   []*post
	superuser bool
}

const (
	maxTitleSize  = 30
	maxAuthorSize = 28
	maxBodySize   = 8192
)

func newPost(title string, author string, body string, tripcoder *Tripcoder) *post {
	var (
		tripcode  string
		superuser bool
	)

	if tripcoder != nil {
		author, tripcode = tripcoder.code(author)
		superuser = tripcoder.isSuperuser(tripcode)
	}

	return &post{
		Title:     title,
		Author:    author,
		Tripcode:  tripcode,
		Body:      body,
		time:      time.Now(),
		uuid:      uuid.New().String(),
		superuser: superuser,
	}
}

func (p *post) IsOriginalPoster() bool {
	return p.Parent != nil && p.Parent.Tripcode != "" && p.Tripcode == p.Parent.Tripcode
}

func (p *post) IsSuperuser() bool {
	return p.superuser
}

func (p *post) Time() string {
	return p.time.UTC().Format(time.RFC3339)
}

func (p *post) URL() string {
	return "/p/" + p.uuid
}

func (p *post) delete() {
	p.Title = ""
	p.Author = ""
	p.Body = "deleted"
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

func (p *post) validate() []string {
	var errors []string

	if len(p.Title) > maxTitleSize {
		errors = append(errors, fmt.Sprintf("Title size cannot exceed %d characters.", maxTitleSize))
	}

	if len(p.Author) > maxAuthorSize {
		errors = append(errors, fmt.Sprintf("Author size cannot exceed %d characters.", maxAuthorSize))
	}

	if len(p.Body) > maxBodySize {
		errors = append(errors, fmt.Sprintf("Body size cannot exceed %d characters.", maxBodySize))
	}

	return errors
}
