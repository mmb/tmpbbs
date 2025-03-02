package tmpbbs

import (
	"net/http"
	"strings"
)

type PostPostHandler struct {
	repliesPerPage int
	postStore      *PostStore
	tripCoder      *TripCoder
}

func NewPostPostHandler(repliesPerPage int, postStore *PostStore, tripCoder *TripCoder) *PostPostHandler {
	return &PostPostHandler{
		repliesPerPage: repliesPerPage,
		postStore:      postStore,
		tripCoder:      tripCoder,
	}
}

func (pph PostPostHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	parentID, err := castID(r.PathValue("parentID"))
	if err != nil {
		http.NotFound(w, r)

		return
	}

	// The body has CRLF line endings which blackfriday doesn't handle well. Convert to CR.
	body := strings.ReplaceAll(r.FormValue("body"), "\r\n", "\n")
	p := newPost(r.FormValue("title"), r.FormValue("author"), body, pph.tripCoder)
	pph.postStore.put(p, parentID)

	repliesLastPage := p.Parent.repliesLastPage(pph.repliesPerPage)
	http.Redirect(w, r, p.Parent.repliesPageURL(repliesLastPage, "replies-end"), http.StatusSeeOther)
}
