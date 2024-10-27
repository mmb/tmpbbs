package tmpbbs

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"strings"
)

type tripCoder struct {
	salt []byte
}

func NewTripCoder(salt string, randReader io.Reader) (*tripCoder, error) {
	tc := tripCoder{}

	if salt != "" {
		tc.salt = []byte(salt)
	} else {
		tc.salt = make([]byte, 16)
		_, err := randReader.Read(tc.salt)
		if err != nil {
			return nil, err
		}
	}

	return &tc, nil
}

func (tc tripCoder) code(s string) (string, string) {
	parts := strings.SplitN(s, "#", 2)
	if len(parts) != 2 {
		return s, ""
	}

	hash := sha256.New()
	hash.Write(tc.salt)
	hash.Write([]byte(s))

	return parts[0], fmt.Sprintf("%.10s", hex.EncodeToString(hash.Sum(nil)))
}
