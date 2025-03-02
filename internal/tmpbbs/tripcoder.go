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
	tripCoder := TripCoder{}

	if salt != "" {
		tripCoder.salt = []byte(salt)
	} else {
		tripCoder.salt = make([]byte, 16)

		_, err := randReader.Read(tripCoder.salt)
		if err != nil {
			return nil, err
		}
	}

	return &tripCoder, nil
}

func (tc TripCoder) code(input string) (string, string) {
	parts := strings.SplitN(input, "#", 2)
	if len(parts) != 2 {
		return input, ""
	}

	if parts[1] == "" {
		return input[:len(input)-1], ""
	}

	hash := sha256.New()
	hash.Write(tc.salt)
	hash.Write([]byte(input))

	return parts[0], fmt.Sprintf("%.10s", hex.EncodeToString(hash.Sum(nil)))
}
