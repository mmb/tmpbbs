package tmpbbs

import (
	"encoding/base32"
	"strings"
)

var (
	crockfordEncoding = base32.NewEncoding( //nolint:gochecknoglobals // immutable
		"0123456789ABCDEFGHJKMNPQRSTVWXYZ").WithPadding(base32.NoPadding)
	crockfordReplacer = strings.NewReplacer( //nolint:gochecknoglobals // constant defined in the spec
		"-", "",
		"I", "1",
		"L", "1",
		"O", "0",
	)
)

func crockfordNormalize(s string) string {
	return crockfordReplacer.Replace(strings.ToUpper(s))
}
