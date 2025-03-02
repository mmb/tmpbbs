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

//go:generate gotext update -lang en-US -out catalog.go github.com/mmb/tmpbbs/internal/tmpbbs

type PostGetHandler struct {
	repliesPerPage      int
	cssURLs             []string
	repliesEnabled      bool
	emojiEnabled        bool
	qrCodesEnabled      bool
	basicEmojiParser    parser
	wrappingEmojiParser parser
	markdownParser      parser
	postStore           *PostStore
}

func NewPostGetHandler(repliesPerPage int, cssURLs []string, repliesEnabled bool, emojiEnabled bool, qrCodesEnabled bool, postStore *PostStore) *PostGetHandler {
	var (
		basicEmojiParser    parser
		wrappingEmojiParser parser
	)

	if emojiEnabled {
		basicEmojiParser = newWrappingEmojiParser(nil)
		wrappingEmojiParser = newWrappingEmojiParser(newEmojiSpanWrapper("emoji"))
	}

	return &PostGetHandler{
		repliesPerPage:      repliesPerPage,
		cssURLs:             cssURLs,
		repliesEnabled:      repliesEnabled,
		emojiEnabled:        emojiEnabled,
		qrCodesEnabled:      qrCodesEnabled,
		basicEmojiParser:    basicEmojiParser,
		wrappingEmojiParser: wrappingEmojiParser,
		markdownParser:      newMarkdownParser(),
		postStore:           postStore,
	}
}

func (pgh PostGetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	printer := message.NewPrinter(message.MatchLanguage(r.Header.Get("Accept-Language"), "en-US"))

	if !pgh.postStore.get(id, func(post *post) {
		displayPost := newDisplayPost(post, printer, pgh.basicEmojiParser, pgh.wrappingEmojiParser, pgh.markdownParser)
		if !displayPost.hasRepliesPage(repliesPage, pgh.repliesPerPage) {
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

func (pgh PostGetHandler) renderPost(displayPost *displayPost, repliesPage int, w io.Writer) error {
	return templates.ExecuteTemplate(w, "index.gohtml", map[string]interface{}{
		"cssURLs":        pgh.cssURLs,
		"emojiEnabled":   pgh.emojiEnabled,
		"qrCodesEnabled": pgh.qrCodesEnabled,
		"repliesEnabled": pgh.repliesEnabled,
		"repliesPerPage": pgh.repliesPerPage,
		"post":           displayPost,
		"repliesPage":    repliesPage,
	})
}
