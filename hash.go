package interview_assignment

import (
	"crypto/sha512"
	"encoding/base64"
)

// HashAndEncode will return the base64 encoded hash of the given password.Hash
func HashAndEncode(pass string) string {
	hasher := sha512.New()

	hasher.Write([]byte(pass))

	hash := hasher.Sum([]byte{})

	return base64.StdEncoding.EncodeToString(hash)
}
