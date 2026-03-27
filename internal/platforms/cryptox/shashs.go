package cryptox

import (
	"crypto/sha1"
	"encoding/hex"
)

func HashIdToken(token string) string {
	sum := sha1.Sum([]byte(token))
	return hex.EncodeToString(sum[:])
}
