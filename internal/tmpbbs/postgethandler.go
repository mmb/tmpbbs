package tmpbbs

import (
	"embed"
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
	template            *template.Template
}

//go:embed template
var templateFS embed.FS

func NewPostGetHandler(repliesPerPage int, cssURLs []string, repliesEnabled bool, emojiEnabled bool,
	qrCodesEnabled bool, postStore *PostStore,
) *PostGetHandler {
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
		template:            template.Must(template.New("template").ParseFS(templateFS, "template/*.gohtml")),
	}
}

func (pgh *PostGetHandler) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.Header().Set("Cache-Control", "no-store")

	postID, err := castID(request.PathValue("id"))
	if err != nil {
		http.NotFound(responseWriter, request)

		return
	}

	repliesPage, err := strconv.Atoi(request.URL.Query().Get("p"))
	if err != nil {
		repliesPage = 1
	}

	printer := message.NewPrinter(message.MatchLanguage(request.Header.Get("Accept-Language"), "en-US"))

	if !pgh.postStore.get(postID, func(post *post) {
		displayPost := newDisplayPost(post, printer, pgh.basicEmojiParser, pgh.wrappingEmojiParser, pgh.markdownParser)
		if !displayPost.hasRepliesPage(repliesPage, pgh.repliesPerPage) {
			http.NotFound(responseWriter, request)

			return
		}

		if err = pgh.renderPost(displayPost, repliesPage, responseWriter); err != nil {
			http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		}
	}) {
		http.NotFound(responseWriter, request)
	}
}

func castID(id string) (int, error) {
	if id == "" {
		return 0, nil
	}

	return strconv.Atoi(id)
}

func (pgh *PostGetHandler) renderPost(displayPost *displayPost, repliesPage int, w io.Writer) error {
	return pgh.template.ExecuteTemplate(w, "index.gohtml", map[string]interface{}{
		"cssURLs":        pgh.cssURLs,
		"emojiEnabled":   pgh.emojiEnabled,
		"qrCodesEnabled": pgh.qrCodesEnabled,
		"repliesEnabled": pgh.repliesEnabled,
		"repliesPerPage": pgh.repliesPerPage,
		"post":           displayPost,
		"repliesPage":    repliesPage,
	})
}
