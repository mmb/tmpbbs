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

func (pph PostPostHandler) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	parentID, err := castID(request.PathValue("parentID"))
	if err != nil {
		http.NotFound(responseWriter, request)

		return
	}

	// The body has CRLF line endings which blackfriday doesn't handle well. Convert to CR.
	body := strings.ReplaceAll(request.FormValue("body"), "\r\n", "\n")
	p := newPost(request.FormValue("title"), request.FormValue("author"), body, pph.tripCoder)
	pph.postStore.put(p, parentID)

	repliesLastPage := p.Parent.repliesLastPage(pph.repliesPerPage)
	http.Redirect(responseWriter, request, p.Parent.repliesPageURL(repliesLastPage, "replies-end"), http.StatusSeeOther)
}
