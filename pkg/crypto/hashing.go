package crypto

import (
	"crypto/sha256"
	"encoding/hex"
)

// HashSHA256 generates a SHA256 hash for a given string
func HashSHA256(data string) string {
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}
