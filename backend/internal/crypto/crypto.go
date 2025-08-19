package crypto

import (
	"crypto/rand"
	"encoding/base64"
)

// GenerateRandomBytes generates a cryptographically secure random byte slice of a given length.
func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// GenerateRandomToken generates a Base64 URL-encoded random string of a given length.
func GenerateRandomToken(length int) (string, error) {
	bytes, err := GenerateRandomBytes(length)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}
