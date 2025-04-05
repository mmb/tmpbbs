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

func (wep wrappingEmojiParser) parse(input string) string {
	// This check shouldn't be necessary because it's checked for nil inside
	// emoji.ParseWithWrapper but text/template complains about a nil pointer
	// if it's not checked here.
	if wep.wrapper == nil {
		return emoji.Sprint(input)
	}

	return emoji.SprintWithWrapper(wep.wrapper, input)
}
