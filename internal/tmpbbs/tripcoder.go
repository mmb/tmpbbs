package tmpbbs

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"strings"
)

type TripCoder struct {
	salt []byte
}

func NewTripCoder(salt string, randReader io.Reader) (*TripCoder, error) {
	tc := TripCoder{}

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

func (tc TripCoder) code(s string) (string, string) {
	parts := strings.SplitN(s, "#", 2)
	if len(parts) != 2 {
		return s, ""
	}
	if parts[1] == "" {
		return s[:len(s)-1], ""
	}

	hash := sha256.New()
	hash.Write(tc.salt)
	hash.Write([]byte(s))

	return parts[0], fmt.Sprintf("%.10s", hex.EncodeToString(hash.Sum(nil)))
}
