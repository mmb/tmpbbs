package tmpbbs

import "strings"

var crockfordReplacer = strings.NewReplacer( //nolint:gochecknoglobals // constant defined in the spec
	"-", "",
	"I", "1",
	"L", "1",
	"O", "0",
	"U", "",
)

func crockfordDecode(s string) string {
	return crockfordReplacer.Replace(strings.ToUpper(s))
}
