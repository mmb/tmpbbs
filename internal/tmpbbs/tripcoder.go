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

const randomSaltLength = 16

func NewTripCoder(salt string, randReader io.Reader) (*TripCoder, error) {
	var saltBytes []byte
	if salt == "" {
		saltBytes = make([]byte, randomSaltLength)

		if _, err := randReader.Read(saltBytes); err != nil {
			return nil, err
		}
	} else {
		saltBytes = []byte(salt)
	}

	return &TripCoder{salt: saltBytes}, nil
}

func (tc TripCoder) code(input string) (string, string) {
	parts := strings.SplitN(input, "#", 2) //nolint:mnd // input has two parts, can't change
	if len(parts) != 2 {                   //nolint:mnd // input has two parts, can't change
		return input, ""
	}

	if parts[1] == "" {
		return input[:len(input)-1], ""
	}

	hash := sha256.New()
	hash.Write(tc.salt)       //nolint:errcheck // can't error
	hash.Write([]byte(input)) //nolint:errcheck // can't error

	return parts[0], fmt.Sprintf("%.10s", hex.EncodeToString(hash.Sum(nil)))
}
