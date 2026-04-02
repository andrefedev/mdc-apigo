package cryptox

import (
	"crypto/sha256"
	"encoding/hex"
)

func HashIdToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
