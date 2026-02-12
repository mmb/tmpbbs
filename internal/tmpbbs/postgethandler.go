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

type postGetHandler struct {
	basicEmojiParser    parser
	markdownParser      parser
	wrappingEmojiParser parser
	postStore           *PostStore
	template            *template.Template
	commit              string
	cssURLs             []string
	repliesPerPage      int
	emojiEnabled        bool
	qrCodesEnabled      bool
	repliesEnabled      bool
}

//go:embed template
var templateFS embed.FS

func newPostGetHandler(postStore *PostStore, commit string, cssURLs []string, emojiEnabled bool, repliesEnabled bool,
	repliesPerPage int, qrCodesEnabled bool,
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
		basicEmojiParser:    basicEmojiParser,
		commit:              commit,
		cssURLs:             cssURLs,
		emojiEnabled:        emojiEnabled,
		markdownParser:      newMarkdownParser(),
		postStore:           postStore,
		qrCodesEnabled:      qrCodesEnabled,
		repliesEnabled:      repliesEnabled,
		repliesPerPage:      repliesPerPage,
		template:            template.Must(template.New("template").ParseFS(templateFS, "template/*.gohtml")),
		wrappingEmojiParser: wrappingEmojiParser,
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

		eTag := fmt.Sprintf(`"%s-%d-%d-%s"`, pgh.commit, post.lastUpdate().UnixNano(), repliesPage, language)
		responseWriter.Header().Set("ETag", eTag)

		if ifNoneMatch := request.Header.Get("If-None-Match"); ifNoneMatch != "" {
			for checkETag := range strings.SplitSeq(ifNoneMatch, ",") {
				if strings.TrimSpace(checkETag) == eTag {
					responseWriter.WriteHeader(http.StatusNotModified)

					return
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
		"commit":         pgh.commit,
		"cssURLs":        pgh.cssURLs,
		"emojiEnabled":   pgh.emojiEnabled,
		"maxAuthorSize":  maxAuthorSize,
		"maxBodySize":    maxBodySize,
		"maxTitleSize":   maxTitleSize,
		"post":           displayPost,
		"qrCodesEnabled": pgh.qrCodesEnabled,
		"repliesEnabled": pgh.repliesEnabled,
		"repliesPage":    repliesPage,
		"repliesPerPage": pgh.repliesPerPage,
	})
}
