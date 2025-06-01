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

// The body is application/x-www-form-urlencoded so the max size is 3 times
// the field limits (if every character were encoded) plus 2 bytes for the &
// separators.
const maxRequestBodyBytes = (maxTitleSize+maxAuthorSize+maxBodySize)*3 + 2

func newPostPostHandler(postStore *PostStore, tripcoder *Tripcoder) *postPostHandler {
	return &postPostHandler{
		postStore: postStore,
		tripcoder: tripcoder,
	}
}

func (pph *postPostHandler) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	parentID := cmp.Or(crockfordNormalize(request.PathValue("parentID")), pph.postStore.rootID)
	if !pph.postStore.hasPost(parentID) {
		http.NotFound(responseWriter, request)

		return
	}

	request.Body = http.MaxBytesReader(responseWriter, request.Body, maxRequestBodyBytes)
	defer request.Body.Close()

	if err := request.ParseForm(); err != nil {
		http.Error(responseWriter, err.Error(), http.StatusBadRequest)

		return
	}

	// The body has CRLF line endings which blackfriday doesn't handle well. Convert to CR.
	body := strings.ReplaceAll(request.FormValue("body"), "\r\n", "\n")

	post := newPost(request.FormValue("title"), request.FormValue("author"), body, pph.tripcoder)
	if errors := post.validate(); len(errors) > 0 {
		http.Error(responseWriter, strings.Join(errors, "\n"), http.StatusBadRequest)

		return
	}

	pph.postStore.put(post, parentID)

	http.Redirect(responseWriter, request, post.Parent.repliesPageURL(1, "replies-start"), http.StatusSeeOther)
}
