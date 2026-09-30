package tmpbbs

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log/slog"
	"strings"
)

// A Tripcoder calculates tripcodes from a salt and user input.
type Tripcoder struct {
	superuserTripcodes map[string]struct{}
	salt               []byte
}

const randomSaltLength = 10 // 10 random bytes (80 bits of entropy) Crockford encodes to a 16 character string

// NewTripcoder returns a new Tripcoder with the passed in salt. If the salt
// is empty 10 random bytes are generated that are Crockford Base 32 encoded to
// a 16 character salt.
func NewTripcoder(salt string, superuserTripcodes []string, randReader io.Reader) (*Tripcoder, error) {
	var saltBytes []byte

	if salt == "" {
		randomBytes := make([]byte, randomSaltLength)

		_, err := io.ReadFull(randReader, randomBytes)
		if err != nil {
			return nil, err
		}

		salt = crockfordEncoding.EncodeToString(randomBytes)
		slog.Info("generated tripcode salt", "salt", salt)
	}

	saltBytes = []byte(salt)

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
		return parts[0], ""
	}

	hash := hmac.New(sha256.New, tc.salt)
	hash.Write([]byte(input)) //nolint:errcheck // can't error

	return parts[0], fmt.Sprintf("%.10s", hex.EncodeToString(hash.Sum(nil)))
}

func (tc Tripcoder) isSuperuser(tripcode string) bool {
	_, found := tc.superuserTripcodes[tripcode]

	return found
}
