package tmpbbs

import (
	"cmp"
	"net/http"
	"strings"
)

type postPostHandler struct {
	postStore *PostStore
	tripcoder *Tripcoder
}

func newPostPostHandler(postStore *PostStore, tripcoder *Tripcoder) *postPostHandler {
	return &postPostHandler{
		postStore: postStore,
		tripcoder: tripcoder,
	}
}

func (pph postPostHandler) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	parentUUID := cmp.Or(crockfordDecode(request.PathValue("parentUUID")), pph.postStore.posts[0].uuid)
	if !pph.postStore.hasPost(parentUUID) {
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

	pph.postStore.put(post, parentUUID)

	http.Redirect(responseWriter, request, post.Parent.repliesPageURL(1, "replies-start"), http.StatusSeeOther)
}
