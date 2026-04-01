package tmpbbs

import (
	"cmp"
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strings"

	"github.com/derekparker/trie"
	"github.com/kyokomi/emoji/v2"
)

type emojiSuggestHandler struct {
	trie *trie.Trie
}

type suggestion struct {
	Pictogram  string `json:"pictogram"`
	Suggestion string `json:"suggestion"`
}

const minQueryLength = 2

func newEmojiSuggestHandler() *emojiSuggestHandler {
	tri := trie.New()
	for shortcode, pictogram := range emoji.CodeMap() {
		tri.Add(shortcode, suggestion{
			Suggestion: strings.Trim(emoji.NormalizeShortCode(shortcode), ":"),
			Pictogram:  pictogram,
		})
	}

	return &emojiSuggestHandler{
		trie: tri,
	}
}

// ServeHTTP serves JSON emoji suggestions for autocomplete.
func (ah *emojiSuggestHandler) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.Header().Set("Content-Type", "application/json")

	query := request.URL.Query().Get("q")
	if len(query) < minQueryLength {
		http.Error(responseWriter, fmt.Sprintf("query must be a minimum of %d characters", minQueryLength),
			http.StatusBadRequest)

		return
	}

	uniqueMap := make(map[suggestion]struct{})

	for _, key := range ah.trie.PrefixSearch(query) {
		value, _ := ah.trie.Find(key)

		metaSuggestion, ok := value.Meta().(suggestion)
		if !ok {
			http.Error(responseWriter, "suggestion type assertion failed", http.StatusInternalServerError)

			return
		}

		uniqueMap[metaSuggestion] = struct{}{}
	}

	result := make([]suggestion, len(uniqueMap))
	i := 0

	for suggestion := range uniqueMap {
		result[i] = suggestion
		i++
	}

	slices.SortFunc(result, func(a suggestion, b suggestion) int {
		return cmp.Compare(a.Suggestion, b.Suggestion)
	})

	err := json.NewEncoder(responseWriter).Encode(result)
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)

		return
	}
}
