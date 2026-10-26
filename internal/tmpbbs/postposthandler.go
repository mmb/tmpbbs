package tmpbbs

import (
	"cmp"
	"errors"
	"net/http"
	"strings"
)

type postPostHandler struct {
	postStore *PostStore
	tripcoder *Tripcoder
}

// Field max sizes are in UTF-16 code units. The maximum size of each is 3 bytes
// times 3 if percent encoded. 20 bytes are added for the field names, equals
// signs and & separators.
const maxRequestBodyBytes = (maxTitleSize+maxAuthorSize+maxBodySize)*9 + 20

func newPostPostHandler(postStore *PostStore, tripcoder *Tripcoder) *postPostHandler {
	return &postPostHandler{
		postStore: postStore,
		tripcoder: tripcoder,
	}
}

// ServeHTTP creates new posts and redirects to the top of the parent's replies.
func (pph *postPostHandler) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	request.Body = http.MaxBytesReader(responseWriter, request.Body, maxRequestBodyBytes)

	err := request.ParseForm()
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusBadRequest)

		return
	}

	post := newPost(request.FormValue("title"), request.FormValue("author"), request.FormValue("body"), pph.tripcoder)
	if validationErrors := post.validate(); len(validationErrors) > 0 {
		http.Error(responseWriter, strings.Join(validationErrors, "\n"), http.StatusBadRequest)

		return
	}

	parentID := cmp.Or(crockfordNormalize(request.PathValue("parentID")), pph.postStore.rootID)
	if errors.Is(pph.postStore.put(post, parentID), errParentNotFound) {
		http.NotFound(responseWriter, request)

		return
	}

	http.Redirect(responseWriter, request, post.Parent.repliesPageURL(1, "replies-start"), http.StatusSeeOther)
}
