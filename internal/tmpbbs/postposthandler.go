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
	if !pph.postStore.hasPost(request.PathValue("parentUUID")) {
		http.NotFound(responseWriter, request)

		return
	}

	// The body has CRLF line endings which blackfriday doesn't handle well. Convert to CR.
	body := strings.ReplaceAll(request.FormValue("body"), "\r\n", "\n")

	post := newPost(request.FormValue("title"), request.FormValue("author"), body, pph.tripcoder)
	if errors := post.validate(); len(errors) > 0 {
		http.Error(responseWriter, strings.Join(errors, "\n"), http.StatusBadRequest)

		return
	}

	pph.postStore.put(post, request.PathValue("parentUUID"))

	repliesLastPage := post.Parent.repliesLastPage(pph.repliesPerPage)
	http.Redirect(responseWriter, request, post.Parent.repliesPageURL(repliesLastPage, "replies-end"), http.StatusSeeOther)
}
