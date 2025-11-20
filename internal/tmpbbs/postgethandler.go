package tmpbbs

import (
	"embed"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/text/message"
)

//go:generate gotext update -lang en-US -out catalog.go github.com/mmb/tmpbbs/internal/tmpbbs

type postGetHandler struct {
	basicEmojiParser    parser
	wrappingEmojiParser parser
	markdownParser      parser
	postStore           *PostStore
	template            *template.Template
	cssURLs             []string
	repliesPerPage      int
	repliesEnabled      bool
	emojiEnabled        bool
	qrCodesEnabled      bool
}

//go:embed template
var templateFS embed.FS

func newPostGetHandler(repliesPerPage int, cssURLs []string, repliesEnabled bool, emojiEnabled bool,
	qrCodesEnabled bool, postStore *PostStore,
) *postGetHandler {
	var (
		basicEmojiParser    parser
		wrappingEmojiParser parser
	)

	if emojiEnabled {
		basicEmojiParser = newWrappingEmojiParser(nil)
		wrappingEmojiParser = newWrappingEmojiParser(newEmojiSpanWrapper("emoji"))
	}

	return &postGetHandler{
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

// ServeHTTP serves rendered post pages.
func (pgh *postGetHandler) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	if !pgh.postStore.get(crockfordNormalize(request.PathValue("id")), func(post *post) {
		language := message.MatchLanguage(request.Header.Get("Accept-Language"), "en-US")
		printer := message.NewPrinter(language)
		displayPost := newDisplayPost(post, printer, pgh.basicEmojiParser, pgh.wrappingEmojiParser, pgh.markdownParser)
		repliesPage, err := strconv.Atoi(request.URL.Query().Get("p"))
		if err != nil {
			repliesPage = 1
		}
		if !displayPost.hasRepliesPage(repliesPage, pgh.repliesPerPage) {
			http.NotFound(responseWriter, request)

			return
		}
		responseWriter.Header().Set("Vary", "Accept-Language")
		if Commit != "" {
			eTag := fmt.Sprintf(`"%s-%d-%d-%s"`, Commit, post.lastUpdate().UnixNano(), repliesPage, language)
			responseWriter.Header().Set("ETag", eTag)
			if ifNoneMatch := request.Header.Get("If-None-Match"); ifNoneMatch != "" {
				for checkETag := range strings.SplitSeq(ifNoneMatch, ",") {
					if strings.TrimSpace(checkETag) == eTag {
						responseWriter.WriteHeader(http.StatusNotModified)

						return
					}
				}
			}
		}
		if err = pgh.renderPost(displayPost, repliesPage, responseWriter); err != nil {
			http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		}
	}) {
		http.NotFound(responseWriter, request)
	}
}

func (pgh *postGetHandler) renderPost(displayPost *displayPost, repliesPage int, w io.Writer) error {
	return pgh.template.ExecuteTemplate(w, "index.gohtml", map[string]any{
		"commit":         Commit,
		"cssURLs":        pgh.cssURLs,
		"emojiEnabled":   pgh.emojiEnabled,
		"qrCodesEnabled": pgh.qrCodesEnabled,
		"repliesEnabled": pgh.repliesEnabled,
		"repliesPerPage": pgh.repliesPerPage,
		"post":           displayPost,
		"repliesPage":    repliesPage,
		"maxTitleSize":   maxTitleSize,
		"maxAuthorSize":  maxAuthorSize,
		"maxBodySize":    maxBodySize,
	})
}
