package tmpbbs

import (
	"embed"
	_ "embed"
	"html/template"
	"io"
	"net/http"
	"strconv"

	"golang.org/x/text/message"
)

type postGetHandler struct {
	repliesPerPage int
	cssURLs        []string
	repliesEnabled bool
	emojiParser    emojiParser
	postStore      *postStore
}

func NewPostGetHandler(repliesPerPage int, cssURLs []string, repliesEnabled bool, emojiParser emojiParser, postStore *postStore) *postGetHandler {
	return &postGetHandler{
		repliesPerPage: repliesPerPage,
		cssURLs:        cssURLs,
		repliesEnabled: repliesEnabled,
		emojiParser:    emojiParser,
		postStore:      postStore,
	}
}

func (pgh postGetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")

	id, err := castID(r.PathValue("id"))
	if err != nil {
		http.NotFound(w, r)
		return
	}

	repliesPage, err := strconv.Atoi(r.URL.Query().Get("p"))
	if err != nil {
		repliesPage = 1
	}

	printer := message.NewPrinter(message.MatchLanguage(r.Header.Get("Accept-Language"), "en"))

	if !pgh.postStore.get(id, func(post *post) {
		displayPost := newDisplayPost(post, printer, pgh.emojiParser, markdownToHTML)
		if !displayPost.HasRepliesPage(repliesPage, pgh.repliesPerPage) {
			http.NotFound(w, r)

			return
		}

		err = pgh.renderPost(displayPost, repliesPage, w)
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

//go:embed template
var templateFS embed.FS

var templates = template.Must(template.New("templates").ParseFS(templateFS, "template/*.gohtml"))

func (pgh postGetHandler) renderPost(displayPost *displayPost, repliesPage int, w io.Writer) error {
	return templates.ExecuteTemplate(w, "index.gohtml", map[string]interface{}{
		"cssURLs":        pgh.cssURLs,
		"repliesEnabled": pgh.repliesEnabled,
		"repliesPerPage": pgh.repliesPerPage,
		"post":           displayPost,
		"repliesPage":    repliesPage,
	})
}
