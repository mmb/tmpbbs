package tmpbbs

import (
	"fmt"
	"strings"
)

type emojiSpanWrapper struct {
	class string
}

func newEmojiSpanWrapper(class string) *emojiSpanWrapper {
	return &emojiSpanWrapper{
		class: class,
	}
}

// Wrap wraps emoji in HTML spans with a class for styling.
func (esw emojiSpanWrapper) Wrap(alias string, code string) string {
	description := strings.ReplaceAll(alias[1:len(alias)-1], "_", " ")

	return fmt.Sprintf(`<span title=%q role="img" aria-label=%q class=%q>%s</span>`, description, description, esw.class,
		code)
}
