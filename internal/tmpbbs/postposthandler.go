package tmpbbs

import (
	"net/http"
	"strings"
)

type PostPostHandler struct {
	postStore      *PostStore
	tripcoder      *Tripcoder
	repliesPerPage int
}

func NewPostPostHandler(repliesPerPage int, postStore *PostStore, tripcoder *Tripcoder) *PostPostHandler {
	return &PostPostHandler{
		repliesPerPage: repliesPerPage,
		postStore:      postStore,
		tripcoder:      tripcoder,
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
	p := newPost(request.FormValue("title"), request.FormValue("author"), body, pph.tripcoder)
	pph.postStore.put(p, parentID)

	repliesLastPage := p.Parent.repliesLastPage(pph.repliesPerPage)
	http.Redirect(responseWriter, request, p.Parent.repliesPageURL(repliesLastPage, "replies-end"), http.StatusSeeOther)
}
