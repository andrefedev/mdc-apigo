package cryptox

import (
	"crypto/rand"
	"encoding/hex"
)

const (
	numChars = "1234567890"
)

func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return nil, err
	}
	return b, nil
}

func GenerateRandomString(l int) (string, error) {
	b, err := generateRandomBytes(l)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func GenerateRandomNumberString(l int) (string, error) {
	b, err := generateRandomBytes(l)
	if err != nil {
		return "", err
	}
	for i := range b {
		b[i] = numChars[int(b[i])%len(numChars)]
	}
	return string(b), nil
}
