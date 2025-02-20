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

func (esw emojiSpanWrapper) Wrap(alias string, code string) string {
	description := strings.ReplaceAll(alias[1:len(alias)-1], "_", " ")

	return fmt.Sprintf(`<span title="%s" role="img" aria-label="%s" class="%s">%s</span>`, description, description, esw.class, code)
}
