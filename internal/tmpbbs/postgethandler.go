package tmpbbs

import (
	_ "embed"
	"io"
	"net/http"
	"strconv"

	"golang.org/x/text/message"
)

type postGetHandler struct {
	title          string
	repliesPerPage int
	cssURLs        []string
	repliesEnabled bool
	emojiEnabled   bool
	postStore      *postStore
}

func NewPostGetHandler(title string, repliesPerPage int, cssURLs []string, repliesEnabled bool, emojiEnabled bool, postStore *postStore) *postGetHandler {
	return &postGetHandler{
		title:          title,
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

	if !pgh.postStore.get(id, func(post *post) {
		printer := message.NewPrinter(message.MatchLanguage(r.Header.Get("Accept-Language"), "en"))

		displayPost := newDisplayPost(post, printer, pgh.emojiEnabled)
		if !displayPost.HasRepliesPage(repliesPage, pgh.repliesPerPage) {
			http.NotFound(w, r)

			return
		}

		err = pgh.renderPost(displayPost, w, repliesPage)
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

func (pgh postGetHandler) renderPost(displayPost *displayPost, w io.Writer, repliesPage int) error {
	return templates.ExecuteTemplate(w, "index.gohtml", map[string]interface{}{
		"title":          pgh.title,
		"repliesPerPage": pgh.repliesPerPage,
		"cssURLs":        pgh.cssURLs,
		"repliesEnabled": pgh.repliesEnabled,
		"post":           displayPost,
		"repliesPage":    repliesPage,
	})
}
