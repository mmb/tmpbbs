package tmpbbs

import (
	"encoding/json"
	"net/http"
	"sort"

	"github.com/derekparker/trie"
	"github.com/kyokomi/emoji/v2"
)

type emojiSuggestHandler struct {
	trie *trie.Trie
}

type suggestion struct {
	Suggestion string `json:"suggestion"`
	Pictogram  string `json:"pictogram"`
}

func newEmojiSuggestHandler() *emojiSuggestHandler {
	tri := trie.New()
	for shortcode, pictogram := range emoji.CodeMap() {
		tri.Add(shortcode, suggestion{
			Suggestion: emoji.NormalizeShortCode(shortcode),
			Pictogram:  pictogram,
		})
	}

	return &emojiSuggestHandler{
		trie: tri,
	}
}

func (ah emojiSuggestHandler) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	var keys []string

	responseWriter.Header().Set("Content-Type", "application/json")

	query := request.URL.Query().Get("q")
	if ah.trie.HasKeysWithPrefix(query) {
		keys = ah.trie.PrefixSearch(query)
	} else {
		if _, err := responseWriter.Write([]byte("[]")); err != nil {
			http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		}

		return
	}

	uniqueMap := make(map[suggestion]struct{})

	for _, key := range keys {
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

	sort.Slice(result, func(i int, j int) bool {
		return result[i].Suggestion[1:len(result[i].Suggestion)-1] < result[j].Suggestion[1:len(result[j].Suggestion)-1]
	})

	err := json.NewEncoder(responseWriter).Encode(result)
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)

		return
	}
}
