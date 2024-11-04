package tmpbbs

import (
	"net/http"
	"strings"
)

type postPostHandler struct {
	repliesPerPage int
	postStore      *postStore
	tripCoder      *tripCoder
}

func NewPostPostHandler(repliesPerPage int, postStore *postStore, tripCoder *tripCoder) *postPostHandler {
	return &postPostHandler{
		repliesPerPage: repliesPerPage,
		postStore:      postStore,
		tripCoder:      tripCoder,
	}
}

func (pph postPostHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	parentID, err := castID(r.PathValue("parentID"))
	if err != nil {
		http.NotFound(w, r)
		return
	}

	// The body has CRLF line endings which blackfriday doesn't handle well. Convert to CR.
	body := strings.Replace(r.FormValue("body"), "\r\n", "\n", -1)
	p := newPost(r.FormValue("title"), r.FormValue("author"), body, pph.tripCoder)
	pph.postStore.put(p, parentID)

	http.Redirect(w, r, newDisplayPost(p.Parent, nil, false).repliesPageEndURL(pph.repliesPerPage, "replies-end"), http.StatusSeeOther)
}
