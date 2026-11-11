package tmpbbs

import (
	"cmp"
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strings"

	"github.com/kyokomi/emoji/v2"
)

type emojiSuggestHandler struct {
	emojiSuggestions []emojiSuggestion
}

type emojiSuggestion struct {
	Pictogram  string `json:"pictogram"`
	shortcode  string
	Suggestion string `json:"suggestion"`
}

const (
	maxResults     = 8
	minQueryLength = 2
)

func newEmojiSuggestHandler() *emojiSuggestHandler {
	handler := &emojiSuggestHandler{
		emojiSuggestions: make([]emojiSuggestion, 0, len(emoji.CodeMap())),
	}
	for shortcode, pictogram := range emoji.CodeMap() {
		handler.emojiSuggestions = append(handler.emojiSuggestions, emojiSuggestion{
			Pictogram:  pictogram,
			shortcode:  strings.Trim(shortcode, ":"),
			Suggestion: strings.Trim(emoji.NormalizeShortCode(shortcode), ":"),
		})
	}

	slices.SortFunc(handler.emojiSuggestions, func(a emojiSuggestion, b emojiSuggestion) int {
		return cmp.Compare(a.shortcode, b.shortcode)
	})

	return handler
}

// ServeHTTP serves JSON emoji suggestions for autocomplete.
func (esh *emojiSuggestHandler) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.Header().Set("Content-Type", "application/json")

	query := request.URL.Query().Get("q")
	if len(query) < minQueryLength {
		http.Error(responseWriter, fmt.Sprintf("query must be a minimum of %d characters", minQueryLength),
			http.StatusBadRequest)

		return
	}

	start, _ := slices.BinarySearchFunc(esh.emojiSuggestions, query, func(suggestion emojiSuggestion, prefix string) int {
		return cmp.Compare(suggestion.shortcode, prefix)
	})

	end := start
	for _, suggestion := range esh.emojiSuggestions[start:] {
		if !strings.HasPrefix(suggestion.shortcode, query) {
			break
		}

		end++
	}

	result := slices.Clone(esh.emojiSuggestions[start:end])
	slices.SortFunc(result, func(a emojiSuggestion, b emojiSuggestion) int {
		return a.compare(b, query)
	})
	result = slices.CompactFunc(result, func(a emojiSuggestion, b emojiSuggestion) bool {
		return a.Suggestion == b.Suggestion
	})

	err := json.NewEncoder(responseWriter).Encode(result[:min(len(result), maxResults)])
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)

		return
	}
}

func (es emojiSuggestion) compare(other emojiSuggestion, query string) int {
	return cmp.Or(
		cmp.Compare(es.queryRank(query), other.queryRank(query)),
		cmp.Compare(len(es.Suggestion), len(other.Suggestion)),
		cmp.Compare(es.Suggestion, other.Suggestion),
	)
}

func (es emojiSuggestion) queryRank(query string) int {
	if strings.HasPrefix(es.Suggestion, query) {
		return 0
	}

	return 1
}
