package tmpbbs

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

type tripCoder struct {
	salt string
}

func NewTripCoder(salt string) *tripCoder {
	if salt == "" {
		return nil
	}

	return &tripCoder{
		salt: salt,
	}
}

func (tc tripCoder) code(s string) string {
	parts := strings.SplitN(s, "#", 2)
	if len(parts) != 2 {
		return s
	}

	hash := sha256.New()
	hash.Write([]byte(tc.salt))
	hash.Write([]byte(s))
	hash.Sum(nil)

	return fmt.Sprintf("%s !%.10s", parts[0], hex.EncodeToString(hash.Sum(nil)))
}
