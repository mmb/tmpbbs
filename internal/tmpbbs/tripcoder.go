package tmpbbs

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"strings"
)

// A Tripcoder calculates tripcodes from a salt and user input.
type Tripcoder struct {
	superuserTripcodes map[string]struct{}
	salt               []byte
}

const randomSaltLength = 16

// NewTripcoder returns a new Tripcoder with the passed in salt. If the salt
// is empty a random 16-byte salt is generated.
func NewTripcoder(salt string, superuserTripcodes []string, randReader io.Reader) (*Tripcoder, error) {
	var saltBytes []byte

	if salt == "" {
		saltBytes = make([]byte, randomSaltLength)

		if _, err := randReader.Read(saltBytes); err != nil {
			return nil, err
		}
	} else {
		saltBytes = []byte(salt)
	}

	tripcoder := &Tripcoder{
		salt:               saltBytes,
		superuserTripcodes: make(map[string]struct{}),
	}
	for _, tripcode := range superuserTripcodes {
		tripcoder.superuserTripcodes[tripcode] = struct{}{}
	}

	return tripcoder, nil
}

func (tc Tripcoder) code(input string) (string, string) {
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

func (tc Tripcoder) isSuperuser(tripcode string) bool {
	_, found := tc.superuserTripcodes[tripcode]

	return found
}
