package tmpbbs

import (
	"container/list"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/oklog/ulid/v2"
)

type post struct {
	time                 time.Time
	Parent               *post
	Replies              *list.List
	postsElement         *list.Element
	parentRepliesElement *list.Element
	Title                string
	Author               string
	Tripcode             string
	Body                 string
	id                   string
	superuser            bool
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
		time:      time.Now(),
		Replies:   list.New(),
		Title:     title,
		Author:    author,
		Tripcode:  tripcode,
		Body:      body,
		id:        ulid.Make().String(),
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
	if p.Parent == nil {
		return "/"
	}

	return "/p/" + p.readableID()
}

// bump moves a post to the top of its parent's replies then bumps each of its
// ancestors all the way to the root post.
func (p *post) bump() {
	if p.Parent == nil {
		return
	}

	p.Parent.Replies.MoveToFront(p.parentRepliesElement)
	p.Parent.bump()
}

func (p *post) delete() {
	p.Title = ""
	p.Author = ""
	p.Body = "deleted"
}

func (p *post) readableID() string {
	return strings.ToLower(fmt.Sprintf("%s-%s-%s-%s-%s-%s-%s", p.id[:4], p.id[4:8], p.id[8:12], p.id[12:16],
		p.id[16:20], p.id[20:24], p.id[24:]))
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
	return max(1, int(math.Ceil(float64(p.Replies.Len())/float64(perPage))))
}

// lastUpdate returns the last update time of the most recent reply or the
// time of the post if there are no replies.
func (p *post) lastUpdate() time.Time {
	if lastReply := p.Replies.Front(); lastReply != nil {
		return lastReply.Value.(*post).lastUpdate() //nolint:errcheck,forcetypeassert // only one type
	}

	return p.time.Round(0)
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
