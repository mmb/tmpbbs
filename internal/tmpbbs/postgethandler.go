package tmpbbs

import (
	_ "embed"
	"io"
	"net/http"
	"strconv"

	"golang.org/x/text/message"
)

type postGetHandler struct {
	repliesPerPage int
	cssURLs        []string
	repliesEnabled bool
	emojiEnabled   bool
	postStore      *postStore
}

func NewPostGetHandler(repliesPerPage int, cssURLs []string, repliesEnabled bool, emojiEnabled bool, postStore *postStore) *postGetHandler {
	return &postGetHandler{
		repliesPerPage: repliesPerPage,
		cssURLs:        cssURLs,
		repliesEnabled: repliesEnabled,
		emojiEnabled:   emojiEnabled,
		postStore:      postStore,
	}
}

func (pgh postGetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id, err := castID(r.PathValue("id"))
	if err != nil {
		http.NotFound(w, r)
		return
	}

	repliesPage, err := strconv.Atoi(r.URL.Query().Get("p"))
	if err != nil {
		repliesPage = 1
	}

	if !pgh.postStore.get(id, func(rootPost *post, post *post) {
		printer := message.NewPrinter(message.MatchLanguage(r.Header.Get("Accept-Language"), "en"))

		rootDisplayPost := newDisplayPost(rootPost, printer, pgh.emojiEnabled)

		displayPost := newDisplayPost(post, printer, pgh.emojiEnabled)
		if !displayPost.HasRepliesPage(repliesPage, pgh.repliesPerPage) {
			http.NotFound(w, r)

			return
		}

		w.Header().Set("Cache-Control", "no-store")
		err = pgh.renderPost(displayPost, rootDisplayPost.DisplayTitle(), repliesPage, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}) {
		http.NotFound(w, r)
	}
}

func castID(id string) (int, error) {
	if id == "" {
		return 0, nil
	}

	return strconv.Atoi(id)
}

func (pgh postGetHandler) renderPost(displayPost *displayPost, pageTitle string, repliesPage int, w io.Writer) error {
	return templates.ExecuteTemplate(w, "index.gohtml", map[string]interface{}{
		"cssURLs":        pgh.cssURLs,
		"repliesEnabled": pgh.repliesEnabled,
		"repliesPerPage": pgh.repliesPerPage,
		"post":           displayPost,
		"title":          pageTitle,
		"repliesPage":    repliesPage,
	})
}
