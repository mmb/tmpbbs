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
	trie := trie.New()
	for shortcode, pictogram := range emoji.CodeMap() {
		trie.Add(shortcode, suggestion{
			Suggestion: emoji.NormalizeShortCode(shortcode),
			Pictogram:  pictogram,
		})
	}

	return &emojiSuggestHandler{
		trie: trie,
	}
}

func (ah emojiSuggestHandler) ServeHTTP(reponseWriter http.ResponseWriter, request *http.Request) {
	var keys []string

	reponseWriter.Header().Set("Content-Type", "application/json")

	query := request.URL.Query().Get("q")
	if ah.trie.HasKeysWithPrefix(query) {
		keys = ah.trie.PrefixSearch(query)
	} else {
		if _, err := reponseWriter.Write([]byte("[]")); err != nil {
			http.Error(reponseWriter, err.Error(), http.StatusInternalServerError)
		}

		return
	}

	uniqueMap := make(map[suggestion]struct{})

	for _, key := range keys {
		value, _ := ah.trie.Find(key)

		metaSuggestion, ok := value.Meta().(suggestion)
		if !ok {
			http.Error(reponseWriter, "suggestion type assertion failed", http.StatusInternalServerError)

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

	err := json.NewEncoder(reponseWriter).Encode(result)
	if err != nil {
		http.Error(reponseWriter, err.Error(), http.StatusInternalServerError)

		return
	}
}
