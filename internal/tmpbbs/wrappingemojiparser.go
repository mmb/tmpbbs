package tmpbbs

import (
	"github.com/kyokomi/emoji/v2"
)

type wrappingEmojiParser struct {
	wrapper *emojiSpanWrapper
}

func newWrappingEmojiParser(wrapper *emojiSpanWrapper) *wrappingEmojiParser {
	return &wrappingEmojiParser{
		wrapper: wrapper,
	}
}

func (wep *wrappingEmojiParser) parse(input string) string {
	return emoji.SprintWithWrapper(wep.wrapper, input)
}
